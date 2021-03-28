package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type MqttMessage struct {
	MessageType int
	Data        string
	TopicSuffix string
	Qos         int
}

func (m *MqttMessage) mqttSendMessage(ch chan MqttMessage) {
	select {
	case ch <- *m:
		// message sent
	default:
		fmt.Printf("MQTT message %v to topic %v could not be sent.\n", m.Data, m.TopicSuffix)
	}
	return
}

type StateData struct {
	AgentStatus   string
	ProgramStatus string
}

func (s *StateData) toJson() string {
	return fmt.Sprintf("{\"agentStatus\": \"%v\", \"programStatus\": \"%v\"}", s.AgentStatus, s.ProgramStatus)
}

const (
	MqttShutdown = iota
	MqttPublish  = iota
	MqttEmpty    = iota
)

func getPasswordAndSubscribeToMqttServer(
	mqttConfig MqttConfig,
	waitGroup *sync.WaitGroup,
	mqttMessages chan MqttMessage,
	processMessages chan ProcessMessage,
) {
	waitGroup.Add(1)
	defer waitGroup.Done()

	password := getMqttPassword(mqttConfig.CustomerId, mqttConfig.DeviceId, mqttConfig.PlatformAddress)
	caPool, err := createCaPool(mqttConfig.CaPath)
	if err != nil {
		return
	}

	subscribeToMqttServer(mqttConfig, password, caPool, mqttMessages, processMessages)
}

func getMqttPassword(customerID string, deviceID string, platformAddress string) string {
	password, err := loadPassword()
	if err != nil {
		// need to register with the platform
		for {
			password, err = registerWithPlatform(customerID, deviceID, platformAddress)
			if err == nil {
				log.Printf("Registered with platform")
				break
			} else { // wait and then try again...
				log.Printf("Failed to register... %+v", err)
				time.Sleep(10 * time.Second)
			}
		}
	} else {
		log.Printf("Loaded platform Password from file")
	}

	return password
}

func loadPassword() (string, error) {
	password, err := ioutil.ReadFile("./MQTT_password")
	if err != nil {
		return "", err
	}
	return string(password), err
}

func registerWithPlatform(customerId string, deviceId string, platformAddress string) (string, error) {
	// register the device
	registerAddress := fmt.Sprintf("%+v/api/devices/register/", platformAddress)
	activateAddress := fmt.Sprintf("%+v/api/devices/activate/", platformAddress)
	var jsonStr = []byte(fmt.Sprintf(`{"customer_id":"%+v", "device_id": "%+v"}`, customerId, deviceId))

	resp, err := http.Post(registerAddress, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", fmt.Errorf("failed to register - %+v", err)
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return "", fmt.Errorf("failure to register - %+v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		return "", fmt.Errorf("failure to register - %+v", body)
	}
	password := dat["password"].(string)

	// confirm activation
	jsonStr = []byte(fmt.Sprintf(`{"device_id":"%+v", "password": "%+v"}`, deviceId, password))
	resp2, err2 := http.Post(activateAddress, "application/json", bytes.NewBuffer(jsonStr))
	if err2 != nil {
		return "", fmt.Errorf("failure to activate - %+v", err2)
	}
	defer resp2.Body.Close()
	if !(resp2.StatusCode >= 200 && resp2.StatusCode < 300) {
		return "", fmt.Errorf("failed to activate due to none 200 response %+v", resp2.StatusCode)
	}

	// write the Password to the Password file
	err = ioutil.WriteFile("./MQTT_password", []byte(password), 0600)
	if err != nil {
		return "", fmt.Errorf("unable to write password file, all is lost")
	}

	return password, nil
}

func subscribeToMqttServer(
	mqttConfig MqttConfig,
	mqttPassword string,
	caPool *x509.CertPool,
	mqttMessages chan MqttMessage,
	processMessages chan ProcessMessage,
) {

	publishTopicPrefix := fmt.Sprintf("devices/out/%v/%v/", mqttConfig.CustomerId, mqttConfig.DeviceId)
	subscribeTopicPrefix := fmt.Sprintf("devices/in/%v/%v/", mqttConfig.CustomerId, mqttConfig.DeviceId)

	commandTopic := subscribeTopicPrefix + "command"
	stateTopic := publishTopicPrefix + "state"

	mqttPayloadOnline := (&StateData{"online", "down"}).toJson()
	mqttPayloadOffline := (&StateData{"offline", "down"}).toJson()

	incomingMessageHandler := func(_ mqtt.Client, msg mqtt.Message) {
		fmt.Println(string(msg.Payload()))
	}

	onConnectHandler := func(client mqtt.Client) {
		log.Printf("Connected to mqtt server %s", mqttConfig.MqttAddress)

		// register to receive commands
		token := client.Subscribe(commandTopic, 2, incomingMessageHandler)
		if token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}

		// announce that I'm online
		token = client.Publish(stateTopic, 2, true, mqttPayloadOnline)
		if token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}

		// send information to process to update its state
		processMessage := ProcessMessage{ProcessPushState, ""}
		processMessage.processSendMessage(processMessages)

	}

	onConnectionLostHandler := func(_ mqtt.Client, err error) {
		log.Printf("Connection lost, reason: %v\n", err)
	}

	log.Printf("Connecting to MQTT Server %s", mqttConfig.MqttAddress)

	opts := mqtt.NewClientOptions().AddBroker(mqttConfig.MqttAddress)
	opts.SetClientID("")
	opts.SetKeepAlive(10)
	opts.SetOnConnectHandler(onConnectHandler)
	opts.SetConnectionLostHandler(onConnectionLostHandler)
	opts.SetUsername(mqttConfig.Username)
	opts.SetPassword(mqttPassword)
	opts.SetMaxReconnectInterval(5 * time.Second)
	opts.SetWill(stateTopic, mqttPayloadOffline, 2, true)
	opts.SetTLSConfig(&tls.Config{RootCAs: caPool})

	// create a client and then clock it until we need to stop
	for {
		c := mqtt.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			// fail to connect, have another go in a bit ... TODO: handle quit here
			fmt.Println("Error connection to MQTT broker: ", token.Error())
			time.Sleep(5 * time.Second)
			continue
		} else {
			clockMqtt(c, publishTopicPrefix, subscribeTopicPrefix, mqttMessages)
			break
		}
	}
}

func clockMqtt(c mqtt.Client, publishTopicPrefix string, subscribeTopicPrefix string, messages chan MqttMessage) {
	defer c.Disconnect(250)

	loop := true
	for loop {
		select {
		case m := <-messages:
			{
				switch m.MessageType {
				case MqttShutdown:
					loop = false
					fmt.Println("Got shutdown message.")
				case MqttPublish:
					c.Publish(publishTopicPrefix+m.TopicSuffix, 2, false, m.Data)
				}
			}
		}
	}

	// unsubscribe and update agent_status
	token := c.Unsubscribe(subscribeTopicPrefix + "command")
	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	mqttPayload := (&StateData{"offline", "down"}).toJson()
	token = c.Publish(publishTopicPrefix+"state", 2, true, mqttPayload)
	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

func createCaPool(caPath string) (*x509.CertPool, error) {
	rootCAs := x509.NewCertPool()

	caCrtPem, err := ioutil.ReadFile(caPath)
	if err != nil {
		return rootCAs, fmt.Errorf("failed to read CA certificate")
	}

	ok := rootCAs.AppendCertsFromPEM(caCrtPem)
	if !ok {
		return rootCAs, fmt.Errorf("failed to parse root certificate")
	}

	return rootCAs, nil
}

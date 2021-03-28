package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

type MqttMessage struct {
	MessageType int
	Data        string
	TopicSuffix string
	Qos         int
}

func (m *MqttMessage) MqttSendMessage(ch chan MqttMessage, timeoutMilSec int) {
	select {
	case ch <- *m:
		// message sent
	case <-time.After(time.Duration(timeoutMilSec) * time.Millisecond):
		fmt.Printf("MQTT message %v to topic %v could not be sent.\n", m.Data, m.TopicSuffix)
	}
}

type StateData struct {
	AgentStatus   string
	ProgramStatus string
}

func (s *StateData) ToJson() string {
	return fmt.Sprintf("{\"agentStatus\": \"%v\", \"programStatus\": \"%v\"}", s.AgentStatus, s.ProgramStatus)
}

type MqttCommDetails struct {
	Address    string
	Username   string
	Password   string
	CustomerId string
	DeviceId   string
	CaPath     string
}

const (
	MqttShutdown = iota
	MqttPublish  = iota
	MqttEmpty    = iota
)

func subscribeToMqttServer(
	mqttCommDetails MqttCommDetails,
	customerId string,
	deviceId string,
	waitGroup *sync.WaitGroup,
	mqttMessages chan MqttMessage,
	processMessages chan ProcessMessage,
) {

	publishTopicPrefix := fmt.Sprintf("devices/out/%v/%v/", customerId, deviceId)
	subscribeTopicPrefix := fmt.Sprintf("devices/in/%v/%v/", customerId, deviceId)

	commandTopic := subscribeTopicPrefix + "command"
	stateTopic := publishTopicPrefix + "state"

	mqttPayloadOnline := (&StateData{"online", "down"}).ToJson()
	mqttPayloadOffline := (&StateData{"offline", "down"}).ToJson()

	incomingMessageHandler := func(_ mqtt.Client, msg mqtt.Message) {
		fmt.Println(string(msg.Payload()))
	}

	onConnectHandler := func(client mqtt.Client) {
		log.Printf("Connected to mqtt server %s", mqttCommDetails.Address)

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

		// send information to process to send its state
		processMessage := ProcessMessage{ProcessPushState, ""}
		processMessage.ProcessSendMessage(processMessages, 500)

	}

	onConnectionLostHandler := func(_ mqtt.Client, err error) {
		log.Printf("Connection lost, reason: %v\n", err)
	}

	waitGroup.Add(1) // for the mqtt
	defer waitGroup.Done()
	log.Printf("Connecting to MQTT Server %s", mqttCommDetails.Address)

	opts := mqtt.NewClientOptions().AddBroker(mqttCommDetails.Address)
	opts.SetClientID("")
	opts.SetKeepAlive(10)
	opts.SetOnConnectHandler(onConnectHandler)
	opts.SetConnectionLostHandler(onConnectionLostHandler)
	opts.SetUsername(mqttCommDetails.Username)
	opts.SetPassword(mqttCommDetails.Password)
	opts.SetMaxReconnectInterval(5 * time.Second)
	opts.SetWill(stateTopic, mqttPayloadOffline, 2, true)
	opts.SetTLSConfig(&tls.Config{RootCAs: createCAPool(mqttCommDetails.CaPath)})

	// create a client and then clock it until we need to stop
	for {
		c := mqtt.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			// fail to connect, have another go in a bit ... TODO: handle quit here
			fmt.Println("Error connection to MQTT broker: ", token.Error())
			time.Sleep(5 * time.Second)
			continue
		} else {
			clockMQTT(c, publishTopicPrefix, subscribeTopicPrefix, mqttMessages)
			break
		}
	}
}

func clockMQTT(c mqtt.Client, publishTopicPrefix string, subscribeTopicPrefix string, messages chan MqttMessage) {
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

	mqttPayload := (&StateData{"offline", "down"}).ToJson()
	token = c.Publish(publishTopicPrefix+"state", 2, true, mqttPayload)
	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

func createCAPool(caPath string) *x509.CertPool {
	caCrtPem, err := ioutil.ReadFile(caPath)
	if err != nil {
		panic("Failed to read CA certificate")
	}

	rootCAs := x509.NewCertPool()
	ok := rootCAs.AppendCertsFromPEM(caCrtPem)
	if !ok {
		panic("Failed to parse root certificate")
	}

	return rootCAs
}

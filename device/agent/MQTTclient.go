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

type MqttCommDetails struct {
	Address    string
	Username   string
	Password   string
	CustomerId string
	DeviceId   string
	CaPath     string
}

type MqttTopics struct {
	StatusTopic  string
	CommandTopic string
}

type MqttMessage struct {
	MessageType int
	Data        string
	Topic       string
	Qos         int
}

const (
	SHUTDOWN = iota
	PUBLISH  = iota
)

func incomingMqttMessageHandler(_ mqtt.Client, msg mqtt.Message) {
	fmt.Println(string(msg.Payload()))
}

func subscribeToMqttServer(
	mqttCommDetails MqttCommDetails,
	mqttTopics MqttTopics,
	waitGroup *sync.WaitGroup,
	messages chan MqttMessage,
) {

	waitGroup.Add(1) // for the mqtt
	defer waitGroup.Done()
	log.Printf("Connecting to MQTT Server %s", mqttCommDetails.Address)

	var mqttConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Printf("Connected to mqtt server %s", mqttCommDetails.Address)

		// register to receive commands
		token := client.Subscribe(mqttTopics.CommandTopic, 2, incomingMqttMessageHandler)
		if token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}

		// announce that I'm online
		token = client.Publish(mqttTopics.StatusTopic, 2, true, "{\"status\": \"online\"}")
		if token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
	}

	var mqttConnLostHandler mqtt.ConnectionLostHandler = func(c mqtt.Client, err error) {
		log.Printf("Connection lost, reason: %v\n", err)
	}

	opts := mqtt.NewClientOptions().AddBroker(mqttCommDetails.Address)
	opts.SetClientID("")
	opts.SetKeepAlive(10)
	opts.SetOnConnectHandler(mqttConnectHandler)
	opts.SetConnectionLostHandler(mqttConnLostHandler)
	opts.SetUsername(mqttCommDetails.Username)
	opts.SetPassword(mqttCommDetails.Password)
	opts.SetMaxReconnectInterval(5 * time.Second)
	opts.SetWill(mqttTopics.StatusTopic, "{\"status\": \"offline\"}", 2, true)

	rootCAs := createCAPool(mqttCommDetails.CaPath)
	opts.SetTLSConfig(&tls.Config{RootCAs: rootCAs})

	// create a client and then clock it until we need to stop
	for {
		c := mqtt.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			// fail to connect, have another go in a bit ... TODO: handle quit here
			fmt.Println("Error connection to MQTT broker: ", token.Error())
			time.Sleep(5 * time.Second)
			continue
		} else {
			clockMQTT(c, mqttTopics, messages)
			break
		}
	}
}

func clockMQTT(c mqtt.Client, mqttTopics MqttTopics, messages chan MqttMessage) {
	defer c.Disconnect(250)

	loop := true
	for loop {
		select {
		case m := <-messages:
			{
				switch m.MessageType {
				case SHUTDOWN:
					loop = false
					fmt.Println("Got shutdown message.")
				case PUBLISH:
					c.Publish(mqttTopics.StatusTopic, 2, false, m.Data)
				}
			}
		}
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

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
	address    string
	username   string
	password   string
	customerID string
	deviceID   string
	caPath     string
}

type MqttMessage struct {
	MessageType int
	data        string
	topic       string
	qos         int
}

const (
	SHUTDOWN = iota
	PUBLISH  = iota
)

func subscribeToMqttServer(mqttCommDetails MqttCommDetails, waitGroup *sync.WaitGroup, messages chan MqttMessage) {
	statusTopic := fmt.Sprintf("devices/out/%v/%v/state", mqttCommDetails.customerID, mqttCommDetails.deviceID)
	//commandTopic := fmt.Sprintf("devices/in/%v/%v/command", mqttCommDetails.customerID, mqttCommDetails.deviceID)

	waitGroup.Add(1) // for the mqtt
	defer waitGroup.Done()
	log.Printf("Connecting to MQTT Server %s", mqttCommDetails.address)

	//var mqttCommandHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//	println(string(msg.Payload()))
	//}

	var mqttConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Printf("Connected to mqtt server %s", mqttCommDetails.address)

		// register to receive commands
		//if token := client.Subscribe(commandTopic, 2, mqttCommandHandler); token.Wait() && token.Error() != nil {
		//	log.Fatal(token.Error())
		//}

		// announce that I'm online
		if token := client.Publish(statusTopic, 2, true, "{\"status\": \"online\"}"); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
	}

	var mqttConnLostHandler mqtt.ConnectionLostHandler = func(c mqtt.Client, err error) {
		log.Printf("Connection lost, reason: %v\n", err)
	}

	opts := mqtt.NewClientOptions().AddBroker(mqttCommDetails.address)
	opts.SetClientID("")
	opts.SetKeepAlive(10)
	opts.SetOnConnectHandler(mqttConnectHandler)
	opts.SetConnectionLostHandler(mqttConnLostHandler)
	opts.SetUsername(mqttCommDetails.username)
	opts.SetPassword(mqttCommDetails.password)
	opts.SetMaxReconnectInterval(5 * time.Second)
	opts.SetWill(statusTopic, "{\"status\": \"offline\"}", 2, true)

	rootCAs := createCAPool(mqttCommDetails.caPath)
	opts.SetTLSConfig(&tls.Config{RootCAs: rootCAs})

	// create a client and then clock it until we need to stop
	for {
		c := mqtt.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			// fail to connect, have another go in a bit ... TODO: handle quit here
			fmt.Println(token.Error())
			time.Sleep(1 * time.Second)
			continue
		} else {
			clockMQTT(c, mqttCommDetails.deviceID, messages)
			break
		}
	}
}

func clockMQTT(c mqtt.Client, statusTopic string, messages chan MqttMessage) {
	defer c.Disconnect(250)

	loop := true
	for loop {
		select {
		case m := <-messages:
			{
				switch m.MessageType {
				case SHUTDOWN:
					loop = false
					print("got shutdown message\n")
				case PUBLISH:
					//c.Publish(fmt.Sprintf("devices/%v/stdout", deviceID), 2, false, m.data)
				}
			}
		}
	}

	// unsubscribe and update status
	//if token := c.Unsubscribe(fmt.Sprintf("devices/%v/command", deviceID)); token.Wait() && token.Error() != nil {
	//	log.Fatal(token.Error())
	//}
	if token := c.Publish(statusTopic, 2, true, "{\"status\": \"offline\"}"); token.Wait() && token.Error() != nil {
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

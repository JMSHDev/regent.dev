package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"sync"
	"time"
)

type MQTTServerDetails struct {
	address  string
	username string
	password string
}

func subscribeToMqttServer(mqttServer MQTTServerDetails, waitGroup *sync.WaitGroup, deviceID string, messages chan string) {
	waitGroup.Add(1)
	defer waitGroup.Done()
	log.Printf("Connecting to MQTT Server %s", mqttServer.address)

	var mqttCommandHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		print(string(msg.Payload()))
	}

	var mqttConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Printf("Connected to mqtt server %s", mqttServer.address)

		// register to receive commands
		if token := client.Subscribe(fmt.Sprintf("devices/%v/command", deviceID), 2, mqttCommandHandler); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}

		// announce that I'm online
		if token := client.Publish(fmt.Sprintf("devices/%v/status", deviceID), 2, true, "Online"); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
	}

	var mqttConnLostHandler mqtt.ConnectionLostHandler = func(c mqtt.Client, err error) {
		log.Printf("Connection lost, reason: %v\n", err)
	}

	opts := mqtt.NewClientOptions().AddBroker(mqttServer.address)
	opts.SetClientID("")
	opts.SetKeepAlive(10)
	opts.SetOnConnectHandler(mqttConnectHandler)
	opts.SetConnectionLostHandler(mqttConnLostHandler)
	opts.SetUsername(mqttServer.username)
	opts.SetPassword(mqttServer.password)
	opts.SetMaxReconnectInterval(5 * time.Second)
	opts.SetWill(fmt.Sprintf("devices/%v/status", deviceID), "Offline", 2, true)

	// create a client and then clock it until we need to stop
	for {
		c := mqtt.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			// fail to connect, have another go in a bit ... TODO: handle quit here
			time.Sleep(1 * time.Second)
			continue
		} else {
			clockMQTT(c, deviceID, messages)
			break
		}
	}

	time.Sleep(1 * time.Second)
}

func clockMQTT(c mqtt.Client, deviceID string, messages chan string) {
	defer c.Disconnect(250)

	loop := true
	for loop {
		select {
		case m := <-messages:
			{
				if m == "shutdown" {
					loop = false
					print("got shutdown message\n")
				}
			}
		}
	}

	// unsubscribe and update status
	if token := c.Unsubscribe(fmt.Sprintf("devices/%v/command", deviceID)); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	if token := c.Publish(fmt.Sprintf("devices/%v/status", deviceID), 2, true, "Offline"); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

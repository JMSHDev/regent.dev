package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type MQTTServerDetails struct {
	address  string
	username string
	password string
}

func LaunchMqttServers(mqttServers []MQTTServerDetails, deviceID string) {
	var waitGroup sync.WaitGroup
	for _, s := range mqttServers {
		waitGroup.Add(1)
		go subscribeToMqttServer(s, &waitGroup, deviceID)
	}
	//waitGroup.Wait()
}

func subscribeToMqttServer(mqttServer MQTTServerDetails, waitGroup *sync.WaitGroup, deviceID string) {
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		// check to see if we need to quit
		select {
		case sig := <-sigs:
			log.Print("Exit signal received")
			log.Print(sig)
			break
		default:
		}
		c := mqtt.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			// fail to connect, have another go in a bit
			time.Sleep(1 * time.Second)
			continue
		} else {
			runUntilExit(c, sigs, deviceID)
			break
		}
	}

	time.Sleep(1 * time.Second)
}

func runUntilExit(c mqtt.Client, sigs chan os.Signal, deviceID string) {
	defer c.Disconnect(250)

	// wait for exit
	sig := <-sigs
	log.Print("Exit signal received")
	log.Print(sig)

	// unsubscribe and update status
	if token := c.Unsubscribe(fmt.Sprintf("devices/%v/command", deviceID)); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	if token := c.Publish(fmt.Sprintf("devices/%v/status", deviceID), 2, true, "Offline"); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

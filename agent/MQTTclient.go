package main

import (
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

func LaunchMqttServers(mqttServers []MQTTServerDetails) {
	var waitGroup sync.WaitGroup
	for _, s := range mqttServers {
		waitGroup.Add(1)
		go subscribeToMqttServer(s, &waitGroup)
	}
	waitGroup.Wait()
}

func subscribeToMqttServer(mqttServer MQTTServerDetails, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	log.Printf("Connecting to MQTT Server %s", mqttServer.address)

	var mqttMessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		print(string(msg.Payload()))
	}

	var mqttConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Printf("Connected to mqtt server %s", mqttServer.address)
		if token := client.Subscribe("IoTEnabler/#", 0, mqttMessageHandler); token.Wait() && token.Error() != nil {
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
			runUntilExit(c, sigs)
			break
		}
	}

	time.Sleep(1 * time.Second)
}

func runUntilExit(c mqtt.Client, sigs chan os.Signal) {
	defer c.Disconnect(250)

	// wait for exit
	sig := <-sigs
	log.Print("Exit signal received")
	log.Print(sig)

	// unsubscribe and disconnect
	if token := c.Unsubscribe("IoTEnabler/#"); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

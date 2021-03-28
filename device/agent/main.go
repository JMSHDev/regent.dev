package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	mqttConfig, processConfig, err := loadConfig()
	if err != nil {
		saveDefaultConfig()
		log.Fatal("Config not found - created default")
	}

	mqttMessages := make(chan MqttMessage, 100)
	processMessages := make(chan ProcessMessage, 100)
	var waitGroup sync.WaitGroup // wait for everything to finish so can safely shutdown

	go getPasswordAndSubscribeToMqttServer(mqttConfig, &waitGroup, mqttMessages, processMessages)
	go launchProcess(processConfig, processMessages, mqttMessages, &waitGroup)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			// check to see if we need to quit
			select {
			case sig := <-sigs:
				log.Print("Exit signal received\n")
				log.Print(sig)
				mqttMessage := MqttMessage{MqttShutdown, "", "", 2}
				mqttMessage.mqttSendMessage(mqttMessages)
				processMessage := ProcessMessage{ProcessShutdown, ""}
				processMessage.processSendMessage(processMessages)
				break
			default:
			}
		}
	}()

	print("Waiting for completion\n")
	waitGroup.Wait()
	print("done\n")
}

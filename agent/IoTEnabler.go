package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		saveDefaultConfig()
		log.Fatal("Config not found - created default")
	}

	log.Printf("%v\n", config)

	mqttServer := MQTTServerDetails{
		address:  "FlamingHellfish:1883",
		username: "",
		password: "",
	}

	mqttMessages := make(chan MQTTMessage)
	processMessages := make(chan string)
	var waitGroup sync.WaitGroup // wait for everything to finish so can safely shutdown

	go subscribeToMqttServer(mqttServer, &waitGroup, config.DeviceID, mqttMessages)
	LaunchProcess(config.PathToExecutable,
		config.Arguments,
		processMessages,
		mqttMessages,
		config.AutoRestart,
		config.RestartDelayMs,
		config.DeviceID,
		&waitGroup)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			// check to see if we need to quit
			select {
			case sig := <-sigs:
				log.Print("Exit signal received\n")
				log.Print(sig)
				mqttMessages <- MQTTMessage{SHUTDOWN, "", "", 2}
				processMessages <- "shutdown"
				break
			default:
			}
		}
	}()

	print("Waiting for completion\n")
	waitGroup.Wait()
	print("done\n")
}

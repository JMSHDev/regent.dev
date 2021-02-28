package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	// http /api/devices/register/
	registerWithPlatform(config.CustomerID, config.DeviceID)

	mqttServer := MQTTServerDetails{
		address:  config.MQTTAddress,
		username: config.MQTTUsername,
		password: config.MQTTPassword,
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

func registerWithPlatform(customerId string, deviceId string) (string, error) {
	// register the device
	var jsonStr = []byte(fmt.Sprintf(`{"customer_id":"%+v", "device_id": "%+v"}`, customerId, deviceId))
	resp, err := http.Post("http://localhost/api/devices/register/", "application/json", bytes.NewBuffer(jsonStr))
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

	fmt.Println(password)

	// confirm activation
	jsonStr = []byte(fmt.Sprintf(`{"device_id":"%+v", "password": "%+v"}`, deviceId, password))
	resp2, err2 := http.Post("http://localhost/api/devices/activate/", "application/json", bytes.NewBuffer(jsonStr))
	if err2 != nil {
		return "", fmt.Errorf("failure to activate - %+v", err2)
	}
	defer resp2.Body.Close()
	if !(resp2.StatusCode >= 200 && resp2.StatusCode < 300) {
		return "", fmt.Errorf("failure to activate - non 200 response when activating")
	}

	// write the password to the password file
	err = ioutil.WriteFile("./MQTT_password", []byte(password), 0600)
	if err != nil {
		return "", fmt.Errorf("Unable to write password file")
	}

	os.Exit(0)
	return password, nil
}

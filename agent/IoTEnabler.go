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
	"time"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		saveDefaultConfig()
		log.Fatal("Config not found - created default")
	}

	log.Printf("%v\n", config)

	// TODO: make this asynchronous
	password := getMqttPassword(config.CustomerID, config.DeviceID)

	mqttServer := MQTTServerDetails{
		address:  config.MQTTAddress,
		username: config.DeviceID,
		password: password,
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

func getMqttPassword(customerID string, deviceID string) string {
	password, err := loadPassword()
	if err != nil {
		// need to register with the platform
		for {
			password, err = registerWithPlatform(customerID, deviceID)
			if err == nil {
				log.Printf("Registered with platform")
				break
			} else { // wait and then try again...
				log.Printf("Failed to register... %+v", err)
				time.Sleep(10 * time.Second)
			}
		}
	} else {
		log.Printf("Loaded platform password from file")
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
	// write the password to the password file
	err = ioutil.WriteFile("./MQTT_password", []byte(password), 0600)
	if err != nil {
		return "", fmt.Errorf("Unable to write password file")
	}

	// confirm activation
	jsonStr = []byte(fmt.Sprintf(`{"device_id":"%+v", "password": "%+v"}`, deviceId, password))
	resp2, err2 := http.Post("http://localhost/api/devices/activate/", "application/json", bytes.NewBuffer(jsonStr))
	if err2 != nil {
		// delete the password file, since activation failed
		e := os.Remove("./MQTT_password")
		if e != nil {
			// bugger! - don't know what to do here
			return "", fmt.Errorf("failed to activate due to %+v but cannot delete mqtt password file?! %+v", err2, e)
		}
		return "", fmt.Errorf("failure to activate - %+v", err2)
	}
	defer resp2.Body.Close()
	if !(resp2.StatusCode >= 200 && resp2.StatusCode < 300) {
		// delete the password file, since activation failed
		e := os.Remove("./MQTT_password")
		if e != nil {
			// bugger! - don't know what to do here
			return "", fmt.Errorf("failed to activate due to none 200 response %+v, but cannot delete mqtt password file?! %+v", resp2.StatusCode, e)
		}

		return "", fmt.Errorf("failed to activate due to none 200 response %+v", resp2.StatusCode)
	}

	return password, nil
}

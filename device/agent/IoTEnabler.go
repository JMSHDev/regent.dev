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
	password := getMqttPassword(config.CustomerId, config.DeviceId, config.PlatformAddress)

	mqttCommDetails := MqttCommDetails{
		Address:    config.MqttAddress,
		Username:   config.DeviceId,
		Password:   password,
		CustomerId: config.CustomerId,
		DeviceId:   config.DeviceId,
		CaPath:     config.CaPath,
	}

	mqttMessages := make(chan MqttMessage)
	processMessages := make(chan string)
	var waitGroup sync.WaitGroup // wait for everything to finish so can safely shutdown

	go subscribeToMqttServer(mqttCommDetails, config.CustomerId, config.DeviceId, &waitGroup, mqttMessages)
	LaunchProcess(config.PathToExecutable,
		config.Arguments,
		processMessages,
		mqttMessages,
		config.AutoRestart,
		config.RestartDelayMs,
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
				mqttMessages <- MqttMessage{SHUTDOWN, "", "", 2}
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

func getMqttPassword(customerID string, deviceID string, platformAddress string) string {
	password, err := loadPassword()
	if err != nil {
		// need to register with the platform
		for {
			password, err = registerWithPlatform(customerID, deviceID, platformAddress)
			if err == nil {
				log.Printf("Registered with platform")
				break
			} else { // wait and then try again...
				log.Printf("Failed to register... %+v", err)
				time.Sleep(10 * time.Second)
			}
		}
	} else {
		log.Printf("Loaded platform Password from file")
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

func registerWithPlatform(customerId string, deviceId string, platformAddress string) (string, error) {
	// register the device
	registerAddress := fmt.Sprintf("%+v/api/devices/register/", platformAddress)
	activateAddress := fmt.Sprintf("%+v/api/devices/activate/", platformAddress)
	var jsonStr = []byte(fmt.Sprintf(`{"customer_id":"%+v", "device_id": "%+v"}`, customerId, deviceId))

	resp, err := http.Post(registerAddress, "application/json", bytes.NewBuffer(jsonStr))
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
	// write the Password to the Password file
	err = ioutil.WriteFile("./MQTT_password", []byte(password), 0600)
	if err != nil {
		return "", fmt.Errorf("Unable to write Password file")
	}

	// confirm activation
	jsonStr = []byte(fmt.Sprintf(`{"device_id":"%+v", "password": "%+v"}`, deviceId, password))
	resp2, err2 := http.Post(activateAddress, "application/json", bytes.NewBuffer(jsonStr))
	if err2 != nil {
		// delete the Password file, since activation failed
		e := os.Remove("./MQTT_password")
		if e != nil {
			// bugger! - don't know what to do here
			return "", fmt.Errorf("failed to activate due to %+v but cannot delete mqtt Password file?! %+v", err2, e)
		}
		return "", fmt.Errorf("failure to activate - %+v", err2)
	}
	defer resp2.Body.Close()
	if !(resp2.StatusCode >= 200 && resp2.StatusCode < 300) {
		// delete the Password file, since activation failed
		e := os.Remove("./MQTT_password")
		if e != nil {
			// bugger! - don't know what to do here
			return "", fmt.Errorf("failed to activate due to none 200 response %+v, but cannot delete mqtt Password file?! %+v", resp2.StatusCode, e)
		}

		return "", fmt.Errorf("failed to activate due to none 200 response %+v", resp2.StatusCode)
	}

	return password, nil
}

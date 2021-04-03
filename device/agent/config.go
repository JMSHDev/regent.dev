package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type Config struct {
	PathToExecutable string
	Arguments        string
	AutoRestart      bool
	RestartDelayMs   int
	DeviceId         string
	CustomerId       string
	MqttAddress      string
	PlatformAddress  string
	CaPath           string
}

type ProcessConfig struct {
	PathToExecutable string
	Arguments        string
	AutoRestart      bool
	RestartDelayMs   int
}

type MqttConfig struct {
	MqttAddress     string
	PlatformAddress string
	CustomerId      string
	DeviceId        string
	CaPath          string
}

func loadConfig() (MqttConfig, ProcessConfig, error) {
	f, err := os.Open("config.json")
	if err != nil {
		//no valid config found - make a new one and return it
		return MqttConfig{}, ProcessConfig{}, err
	}
	defer f.Close()

	data, _ := ioutil.ReadAll(f)

	var config Config
	jsonErr := json.Unmarshal(data, &config)
	if jsonErr != nil {
		log.Printf("%v\n", jsonErr)
	}

	mqttConfig := MqttConfig{
		MqttAddress:     config.MqttAddress,
		PlatformAddress: config.PlatformAddress,
		CustomerId:      config.CustomerId,
		DeviceId:        getDeviceId(),
		CaPath:          config.CaPath,
	}

	processConfig := ProcessConfig{
		PathToExecutable: config.PathToExecutable,
		Arguments:        config.Arguments,
		AutoRestart:      config.AutoRestart,
		RestartDelayMs:   config.RestartDelayMs,
	}

	return mqttConfig, processConfig, nil
}

func getDeviceId() string {
	macAddresses, err := getMacAddress()
	if err != nil {
		log.Fatal(err)
	}
	return macAddresses[0]
}

func getMacAddress() ([]string, error) {
	networkInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	if len(networkInterfaces) == 0 {
		return nil, fmt.Errorf("now valid mac address")
	}

	var macAddresses []string
	for _, networkInterface := range networkInterfaces {
		macAddress := networkInterface.HardwareAddr.String()
		if macAddress != "" {
			macAddresses = append(macAddresses, macAddress)
		}
	}
	return macAddresses, nil
}

func saveDefaultConfig() Config {
	f, err := os.Create("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	defaultConfig := Config{
		PathToExecutable: "./ExampleApp",
		Arguments:        "",
		AutoRestart:      true,
		RestartDelayMs:   10000,
		DeviceId:         getDeviceId(),
		CustomerId:       "sample_id",
		MqttAddress:      "ssl://localhost:8883",
		PlatformAddress:  "http://localhost",
		CaPath:           "./ca.crt",
	}
	jsonValue, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	f.Write(jsonValue)
	return defaultConfig
}

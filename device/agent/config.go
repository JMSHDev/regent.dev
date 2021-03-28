package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	Username        string
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
		Username:        config.DeviceId,
		CustomerId:      config.CustomerId,
		DeviceId:        config.DeviceId,
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
		DeviceId:         "deviceID",
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

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
	DeviceID         string
	CustomerID       string
	MQTTAddress      string
	PlatformAddress  string
	// MQTTUsername     string
	// MQTTPassword     string
}

func loadConfig() (Config, error) {
	f, err := os.Open("config.json")
	if err != nil {
		//no valid config found - make a new one and return it
		return Config{}, err
	}
	defer f.Close()

	data, _ := ioutil.ReadAll(f)

	var config Config
	jsonErr := json.Unmarshal(data, &config)
	if jsonErr != nil {
		log.Printf("%v\n", jsonErr)
	}
	return config, nil
}

func saveDefaultConfig() Config {
	f, err := os.Create("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	defaultConfig := Config{
		PathToExecutable: "./program.exe",
		Arguments:        "",
		AutoRestart:      true,
		RestartDelayMs:   10000,
		DeviceID:         "deviceID",
		CustomerID:       "sample_id",
		MQTTAddress:      "localhost:1883",
		PlatformAddress:  "http://localhost",
		//MQTTUsername:     "", // disabled since we want to be using the
		//MQTTPassword:     "",
	}
	jsonValue, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	f.Write(jsonValue)
	return defaultConfig
}

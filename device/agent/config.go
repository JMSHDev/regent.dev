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
		DeviceId:         "DeviceId",
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

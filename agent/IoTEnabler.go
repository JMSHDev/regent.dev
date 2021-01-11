package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Config struct {
	PathToExecutable string
	Arguments        string
}

func main() {
	config, err := loadConfig()
	if err != nil {
		saveDefaultConfig()
		log.Fatal("Config not found - created default")
	}

	log.Printf("%v\n", config)

	cmd := exec.Command(config.PathToExecutable, config.Arguments)
	_, err = cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q\n", out.String())
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
		PathToExecutable: "fish",
		Arguments:        "face",
	}
	jsonValue, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	f.Write(jsonValue)
	return defaultConfig
}

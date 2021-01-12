package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		saveDefaultConfig()
		log.Fatal("Config not found - created default")
	}

	log.Printf("%v\n", config)

	for {
		LaunchProcess(config.PathToExecutable, config.Arguments)
		if !config.AutoRestart {
			fmt.Printf("Process completed...\n")
			break
		} else {
			fmt.Printf("Process exited. Auto restarting\n")
			time.Sleep(time.Duration(config.RestartDelayMs) * time.Millisecond)
		}
	}
}

func LaunchProcess(pathToExecutable string, arguments string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, pathToExecutable, arguments)

	_, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		// was unable to run the program... probably should log & try again after a few seconds
		log.Fatal(err)
	}
	fmt.Printf("%q\n", out.String())
}

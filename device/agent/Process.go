package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"
)

func LaunchProcess(
	pathToExecutable string,
	arguments string,
	inputMessages chan string,
	mqttMessages chan MqttMessage,
	autoRestart bool,
	restartDelayMs int,
	deviceID string,
	waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for {
			exitNow := launchProcessAux(pathToExecutable, arguments, inputMessages, mqttMessages, deviceID)
			if exitNow {
				break
			}
			if !autoRestart {
				fmt.Printf("Process completed...\n")
				break
			} else {
				fmt.Printf("Process exited. Auto restarting\n")
				time.Sleep(time.Duration(restartDelayMs) * time.Millisecond)
			}
		}
	}()
}

func launchProcessAux(pathToExecutable string, arguments string, inputMessages chan string, mqttMessages chan MqttMessage, deviceID string) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, pathToExecutable, arguments)

	_, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan error)
	go func() {
		//funcmqttMessages <- MqttMessage{PUBLISH, "Process started at: " + time.Now().String(), fmt.Sprintf("devices/%v/start", deviceID), 2}
		runResult := cmd.Run()
		//mqttMessages <- MqttMessage{PUBLISH, "Process stopped at: " + time.Now().String(), fmt.Sprintf("devices/%v/stop", deviceID), 2}
		ch <- runResult
	}()

	buf := bufio.NewReader(out) // Notice that this is not in a loop
	var currentLine []byte

	for {
		select {
		case err = <-ch:
			return false
		case command := <-inputMessages:
			if command == "shutdown" {
				print("shutdown process now\n")
				return true
			}
		default:
		}

		bytes, err := buf.ReadBytes('\n')
		if err == io.EOF {
			currentLine = append(currentLine, bytes...)
		} else if err != nil {
			break // some other nasty error
		} else {
			currentLine = append(currentLine, bytes...)
			print(string(currentLine))
			//mqttMessages <- MqttMessage{PUBLISH, string(currentLine), fmt.Sprintf("devices/%v/stdout", deviceID), 2}
			currentLine = []byte{}
		}
	}
	return false
}

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
	waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for {
			exitNow := launchProcessAux(pathToExecutable, arguments, inputMessages, mqttMessages)
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

func launchProcessAux(
	pathToExecutable string,
	arguments string,
	inputMessages chan string,
	mqttMessages chan MqttMessage,
) bool {

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
		messageStart := "Process started at: " + time.Now().String()
		mqttMessages <- MqttMessage{PUBLISH, messageStart, "start", 2}
		runResult := cmd.Run()
		messageStop := "Process stopped at: " + time.Now().String()
		mqttMessages <- MqttMessage{PUBLISH, messageStop, "stop", 2}
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
				fmt.Println("Shutdown process now.")
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
			mqttMessages <- MqttMessage{PUBLISH, string(currentLine), "stdout", 2}
			currentLine = []byte{}
		}
	}
	return false
}

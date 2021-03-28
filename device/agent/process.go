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

type ProcessConfig struct {
	PathToExecutable string
	Arguments        string
	AutoRestart      bool
	RestartDelayMs   int
}

type ProcessMessage struct {
	MessageType int
	Message     string
}

func (m *ProcessMessage) ProcessSendMessage(ch chan ProcessMessage, timeoutMilSec int) {
	select {
	case ch <- *m:
		// message sent
	case <-time.After(time.Duration(timeoutMilSec) * time.Millisecond):
		fmt.Printf("Process message %v of type %v could not be sent.\n", m.Message, m.MessageType)
	}
}

func ProcessReadMessage(ch chan ProcessMessage, timeoutMilSec int) ProcessMessage {
	select {
	case m := <-ch:
		// message read
		return m
	case <-time.After(time.Duration(timeoutMilSec) * time.Millisecond):
		fmt.Printf("No message received from channel.\n")
		return ProcessMessage{ProcessEmpty, ""}
	}
}

const (
	ProcessShutdown  = iota
	ProcessPushState = iota
	ProcessEmpty     = iota
)

func LaunchProcess(
	processConfig ProcessConfig,
	processMessages chan ProcessMessage,
	mqttMessages chan MqttMessage,
	waitGroup *sync.WaitGroup,
) {

	waitGroup.Add(1)
	defer waitGroup.Done()

	for {
		exitNow := launchProcessAux(
			processConfig.PathToExecutable,
			processConfig.Arguments,
			processMessages,
			mqttMessages,
		)
		if exitNow {
			break
		}
		if !processConfig.AutoRestart {
			fmt.Printf("Process completed...\n")
			break
		} else {
			fmt.Printf("Process exited. Auto restarting\n")
			time.Sleep(time.Duration(processConfig.RestartDelayMs) * time.Millisecond)
		}
	}
}

func launchProcessAux(
	pathToExecutable string,
	arguments string,
	processMessages chan ProcessMessage,
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
		runResult := cmd.Run()

		stateMessageStop := MqttMessage{
			MqttPublish,
			(&StateData{"online", "down"}).ToJson(),
			"state",
			2}
		stateMessageStop.MqttSendMessage(mqttMessages, 2000)

		supervisorMessageStop := MqttMessage{
			MqttPublish,
			"Process stopped at: " + time.Now().String(),
			"supervisor",
			2}
		supervisorMessageStop.MqttSendMessage(mqttMessages, 2000)

		ch <- runResult
	}()

	buf := bufio.NewReader(out) // Notice that this is not in a loop
	var currentLine []byte

	for {
		select {
		case <-ch:
			return false
		case processMessage := <-processMessages:
			if processMessage.MessageType == ProcessShutdown {
				fmt.Println("Shutdown process now.")
				return true
			}
			if processMessage.MessageType == ProcessPushState {
				stateMessageStart := MqttMessage{
					MqttPublish,
					(&StateData{"online", "up"}).ToJson(),
					"state",
					2}
				stateMessageStart.MqttSendMessage(mqttMessages, 2000)

				supervisorMessageStart := MqttMessage{
					MqttPublish,
					"Process started at: " + time.Now().String(),
					"supervisor",
					2}
				supervisorMessageStart.MqttSendMessage(mqttMessages, 2000)
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
			supervisorMessage := MqttMessage{
				MqttPublish,
				string(currentLine),
				"supervisor",
				2}
			supervisorMessage.MqttSendMessage(mqttMessages, 2000)
			currentLine = []byte{}
		}
	}
	return false
}

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

type ProcessMessage struct {
	MessageType int
	Message     string
}

func (m *ProcessMessage) processSendMessage(ch chan ProcessMessage) {
	select {
	case ch <- *m:
		// message sent
	default:
	}
	return
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

	errorChannel := make(chan error)
	go func() {
		stateMessageStart := MqttMessage{
			MqttPublish,
			(&StateData{"online", "up"}).toJson(),
			"state",
			2}
		stateMessageStart.mqttSendMessage(mqttMessages)

		supervisorMessageStart := MqttMessage{
			MqttPublish,
			"Process started at: " + time.Now().String(),
			"supervisor",
			2}
		supervisorMessageStart.mqttSendMessage(mqttMessages)
		runResult := cmd.Run()

		stateMessageStop := MqttMessage{
			MqttPublish,
			(&StateData{"online", "down"}).toJson(),
			"state",
			2}
		stateMessageStop.mqttSendMessage(mqttMessages)

		supervisorMessageStop := MqttMessage{
			MqttPublish,
			"Process stopped at: " + time.Now().String(),
			"supervisor",
			2}
		supervisorMessageStop.mqttSendMessage(mqttMessages)

		errorChannel <- runResult
	}()

	buf := bufio.NewReader(out) // Notice that this is not in a loop
	var currentLine []byte

	for {
		select {
		case <-errorChannel:
			return false
		case processMessage := <-processMessages:
			if processMessage.MessageType == ProcessShutdown {
				fmt.Println("Shutdown process now.")
				return true
			}
			if processMessage.MessageType == ProcessPushState {
				stateMessageStart := MqttMessage{
					MqttPublish,
					(&StateData{"online", "up"}).toJson(),
					"state",
					2}
				stateMessageStart.mqttSendMessage(mqttMessages)

				supervisorMessageStart := MqttMessage{
					MqttPublish,
					"Process started at: " + time.Now().String(),
					"supervisor",
					2}
				supervisorMessageStart.mqttSendMessage(mqttMessages)
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
			supervisorMessage.mqttSendMessage(mqttMessages)
			currentLine = []byte{}
		}
	}
	return false
}

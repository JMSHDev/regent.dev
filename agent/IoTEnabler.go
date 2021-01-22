package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		saveDefaultConfig()
		log.Fatal("Config not found - created default")
	}

	log.Printf("%v\n", config)

	mqttServer := MQTTServerDetails{
		address:  "localhost:1883",
		username: "",
		password: "",
	}

	mqttMessages := make(chan string)
	processMessages := make(chan string)
	var waitGroup sync.WaitGroup // wait for everything to finish so can safely shutdown

	go subscribeToMqttServer(mqttServer, &waitGroup, config.DeviceID, mqttMessages)
	go launchProcess(config.PathToExecutable, config.Arguments, processMessages, config.AutoRestart, config.RestartDelayMs, &waitGroup)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			// check to see if we need to quit
			select {
			case sig := <-sigs:
				log.Print("Exit signal received\n")
				log.Print(sig)
				mqttMessages <- "shutdown"
				processMessages <- "shutdown"
				break
			default:
			}
		}
	}()

	print("Waiting for completion\n")
	waitGroup.Wait()
	print("done\n")
}

func launchProcess(pathToExecutable string, arguments string, messages chan string, autoRestart bool, restartDelayMs int, waitGroup *sync.WaitGroup) {
	go func() {
		waitGroup.Add(1)
		defer waitGroup.Done()

		for {
			exitNow := launchProcessAux(pathToExecutable, arguments, messages)
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

func launchProcessAux(pathToExecutable string, arguments string, messages chan string) bool {
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
		ch <- cmd.Run()
	}()

	buf := bufio.NewReader(out) // Notice that this is not in a loop
	var currentLine []byte

	for {
		select {
		case err = <-ch:
			return false
		default:
		}

		select {
		case command := <-messages:
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
			currentLine = []byte{}
		}
	}
	return false
}

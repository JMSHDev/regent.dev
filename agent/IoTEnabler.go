package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
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

	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan error)
	go func() {
		ch <- cmd.Run()
	}()

	//err = cmd.Start()
	//if err != nil {
	//	// was unable to run the program... probably should log & try again after a few seconds
	//	log.Fatal(err)
	//}

	buf := bufio.NewReader(out) // Notice that this is not in a loop
	var currentLine []byte

	for {
		select {
		case err = <-ch:
			return
		default:
		}

		bytes, err := buf.ReadBytes('\n')
		if err == io.EOF {
			currentLine = append(currentLine, bytes...)
		} else if err != nil {
			break // some othr nasty error
		} else {
			currentLine = append(currentLine, bytes...)
			print(string(currentLine))
			currentLine = []byte{}
		}

	}
}

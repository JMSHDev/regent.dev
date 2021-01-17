package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	argsWithoutProg := os.Args[1:]
	fmt.Printf("%+v\n", argsWithoutProg)

	for i := 1; i < 10; i++ {
		fmt.Println("Zzzzzz")
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Done")
}

package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	argsWithoutProg := os.Args[1:]
	fmt.Printf("%+v\n", argsWithoutProg)
	fmt.Println("Zzzzzz")
	time.Sleep(5 * time.Second)
	fmt.Println("Done")
}

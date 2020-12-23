package main

import (
	"fmt"
	"net"
	"strconv"
)

type udpServer struct {
	name string
	port int
}

type httpClient struct {
	name string
	url  string
}

func main() {
	serverConfig := udpServer{"UDP1", 8081}

	CONNECT := "0.0.0.0:" + strconv.Itoa(serverConfig.port)

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer c.Close()

	for {
		buffer := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Reply: %s\n", string(buffer[0:n]))
	}
}

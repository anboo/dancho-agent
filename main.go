package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type Rotation struct {
	Name string
	Duration int
	Memory int
	Origin float64
	StartTime int
	EndTime int
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		fmt.Printf("Received %s", temp)

		var rotations []Rotation
		json.Unmarshal([]byte(temp), &rotations)

		c.Write([]byte("OK"))
	}

	c.Close()
}

func main() {
	PORT := ":8000"
	l, err := net.Listen("tcp4", PORT); if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept(); if err != nil {
			fmt.Println(err)
			return
		}

		go handleConnection(c)
	}
}


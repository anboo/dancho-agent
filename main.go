package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

type Rotation struct {
	Name string `json:"name"`
	Duration int `json:"duration"`
	Memory int `json:"memory"`
	Origin float64 `json:"origin"`
	StartTime int `json:"start_time"`
	EndTime int `json:"end_time"`
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
		unErr := json.Unmarshal([]byte(temp), &rotations); if unErr != nil {
			log.Fatal(unErr)
		}

		_, errW := c.Write([]byte("OK")); if err != nil {
			log.Fatal(errW)
		}
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


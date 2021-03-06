package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"math/rand"
	"net"
	"net/http"
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

var httpClient *http.Client
var log = logging.MustGetLogger("example")

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	defer c.Close()

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Error(err.Error())
			return
		}

		received := strings.TrimSpace(string(netData))
		bReceived := []byte(received)
		log.Info("Received %s", received)

		var rotations []Rotation
		err = json.Unmarshal(bReceived, &rotations); if err != nil {
			log.Error(err)
			return
		}

		req, err := http.NewRequest("POST", "http://127.0.0.1:8000/api/v1/rotations", bytes.NewBuffer(bReceived)); if err != nil {
			log.Error("Error create HTTP request:" + err.Error())
			return
		}

		res, err := httpClient.Do(req); if err != nil {
			log.Error("Error do request: " + err.Error())
			return
		}

		if res != nil {
			err = res.Body.Close(); if err != nil {
				log.Error(err.Error())
				return
			}
		} else {
			log.Error("Error http request")
			return
		}
	}
}

func main() {
	PORT := ":7731"
	l, err := net.Listen("tcp4", PORT); if err != nil {
		fmt.Println(err)
		return
	}

	httpClient = &http.Client{}

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


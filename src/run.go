package main

import (
	"bufio"
	"fmt"
	"os"
	s "strings"

	"github.com/r3labs/sse"
)

var println = fmt.Println

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func parseCreds(lines []string) map[string]string {
	settings := make(map[string]string)
	// results := make([]string, 2)
	for _, line := range lines {
		// var res []string
		if s.HasPrefix(line, "#") {
			// ignore it
		} else {
			res := s.Split(line, "=")
			settings[res[0]] = res[1]
		}
	}
	return settings
}

func main() {
	// Where the config file is
	credPath := "/Users/eat_sleep_live_skateboarding/Code/go/credentials.txt"

	// The SSE url
	sseURL := "https://api.particle.io/v1/devices/events?access_token="

	// Where we will store our sensor info
	//var sensors []string

	// parse values from config
	_map, _ := readLines(credPath)
	settings := parseCreds(_map)
	println(settings)

	// SSE begins here
	sseURL = sseURL + settings["api-key"]
	println(sseURL)

	client := sse.NewClient(sseURL)
	client.Subscribe("messages", func(msg *sse.Event) {
		// Got some data!
		println(string(msg.Event))
		println(string(msg.Data))
	})

	// events := make(chan *sse.Event)
	//
	// client := sse.NewClient(sseURL)
	// client.SubscribeChan("messages", events)

	// for {
	// 	println("checking")
	// 	event := <-events
	// 	println(event)
	// }
}

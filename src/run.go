package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	s "strings"

	"github.com/r3labs/sse"
)

// Probably not best practice but is conventient
// var println = fmt.Println

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

//EventsReponse is the response of the Particle events API
type EventsResponse struct {
	Data        string `json:"data"`
	PublishedAt string `json:"published_at"`
	CoreID      string `json:"coreid"`
}

//InfluxWriteString
type InfluxWriteString struct {
	// Must conform to:
	// weather,location=us-midwest,season=summer temperature=82 1465839830100400200
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
	fmt.Println(settings)

	// SSE begins here
	sseURL = sseURL + settings["api-key"]
	fmt.Println(sseURL)

	client := sse.NewClient(sseURL)
	go client.Subscribe("messages", func(msg *sse.Event) {
		str := string(msg.Data)
		res := EventsResponse{}
		json.Unmarshal([]byte(str), &res)
		fmt.Println(string(msg.Event))
		fmt.Println(reflect.TypeOf(string(msg.Event)))
		// fmt.Println(reflect.TypeOf(msg.Event))
		// fmt.Println(res)
	})

	var input string
	fmt.Scanln(&input)
}

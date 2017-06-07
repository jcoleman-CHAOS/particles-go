package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/r3labs/sse"
)

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
		if strings.HasPrefix(line, "#") {
			// ignore it
		} else {
			res := strings.Split(line, "=")
			settings[res[0]] = res[1]
		}
	}
	return settings
}

//InfluxWriteString values.
type InfluxWriteString struct {
	// Must conform to:
	// weather,location=us-midwest,season=summer temperature=82 1465839830100400200
}

//GenericSensor struct.
type GenericSensor struct {
	label string
	phen  string //Phenomanon
	unit  string

	location   string  // Qualitatively
	lat        float64 // as float
	long       float64
	Experiment string

	Firmware    string
	PublishRate int64 //milliseconds
}

// AddToEventMap parses Particle API json
func combineEventAndData(se string, sd string) map[string]interface{} {
	m := make(map[string]interface{})
	m["event"] = se
	err := json.Unmarshal([]byte(sd), &m)
	if err != nil {
		if err.Error() == "unexpected end of JSON input" {
			// pass
		} else {
			panic(err)
		}
	}
	return m
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

	// Create the channel to store SSEresponses
	SSEresp := make(chan string, 2)
	counter := 0
	SSEchanIsReady := make(chan bool)

	// SSEres := make(map[string]interface{})
	client := sse.NewClient(sseURL)
	go client.Subscribe("messages", func(msg *sse.Event) {
		if msg.Event != nil {
			SSEresp <- string(msg.Event)
			counter = 0
		} else if msg.Data != nil {
			SSEresp <- string(msg.Data)
			counter = 1
			SSEchanIsReady <- true
		}
	})

	go func() {
		for {
			if <-SSEchanIsReady {
				fmt.Println("\n***")
				res := combineEventAndData(<-SSEresp, <-SSEresp)
				for k, v := range res {
					fmt.Printf("%s: %s\n", k, v)
				}
			}
		}
	}()

	var input string
	fmt.Scanln(&input)
}

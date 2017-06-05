package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/r3labs/sse"
)

// Probably not best practice but is conventient
// var println = fmt.Println

type JSONableSlice []uint8

func (u JSONableSlice) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}

type Test struct {
	Name  string
	Array JSONableSlice
}

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
		if strings.HasPrefix(line, "#") {
			// ignore it
		} else {
			res := strings.Split(line, "=")
			settings[res[0]] = res[1]
		}
	}
	return settings
}

//EventsReponse is the response of the Particle events API.
type EventsResponse struct {
	Data        string `json:"data"`
	PublishedAt string `json:"published_at"`
	CoreID      string `json:"coreid"`
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

// parses Particle API json
func EventsAPIJSON(u []uint8) map[string]interface{} {
	s := string(u)
	var formattedJSON map[string]interface{}
	err := json.Unmarshal([]byte(s), &formattedJSON)
	if err != nil {
		if err.Error() == "unexpected end of JSON input" {
			// pass
		} else {
			panic(err)
		}
	}
	return formattedJSON
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
	client.Subscribe("messages", func(msg *sse.Event) {
		fmt.Println("***")
		fmt.Println("raw")
		fmt.Println(string(msg.Event))
		fmt.Println(string(msg.Data))
		APIres := EventsAPIJSON(msg.Data)
		if len(APIres) == 0 {
			fmt.Println("if ZERO")
			// do nothing
		} else if len(APIres) < 3 {
			fmt.Println("if < THREE")
			for k, v := range APIres {
				fmt.Printf("key[%s] value[%s] type:%s\n", k, v, reflect.TypeOf(v))
			}
		} else {
			fmt.Println("ELSE")
			fmt.Println(string(msg.Event))
			for k, v := range APIres {
				fmt.Printf("key[%s] value[%s] type:%s\n", k, v, reflect.TypeOf(v))
			}
			// APIres["event"] = string(msg.Event)
			fmt.Println("")
		}
	})

	var input string
	fmt.Scanln(&input)
}

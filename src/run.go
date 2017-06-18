package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/r3labs/sse"
)

const (
	value = "value"
)

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

// AllParticlesMeta is an example of what a meta data obj might look like
type AllParticlesMeta struct {
	numParticles int
	numConnected int
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

func urlResp(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Body error" + err.Error())
	}
	return body
}

func iterMap(m map[string]interface{}) {
	for k, v := range m {
		fmt.Printf("\n%s: %v", k, v)
		// fmt.Printf("\n%s is type %T: %v is type %T", k, k, v, v)
	}
	fmt.Println("\n ")
}

func allParticlesCurl(token string) []byte {
	// This is the device URL, contains info on all particles
	devicesAPI := "https://api.particle.io/v1/devices/?access_token=" + token
	fmt.Printf("Checking: %s\n", devicesAPI)
	resp := urlResp(devicesAPI)
	return resp
}

// AddToEventMap parses Particle API json
func combineEventAndData(se string, sd string) map[string]interface{} {
	m := make(map[string]interface{})
	m["event"] = se
	sortEvent(m["event"].(string))
	err := json.Unmarshal([]byte(sd), &m)
	if err != nil {
		if err.Error() == "unexpected end of JSON input" {
			// pass
		} else {
			fmt.Println("It died trying to CombineEventAndData responses")
			panic(err)
		}
	}
	return m
}

func sortEvent(event string) {
	// fmt.Println(event)
	switch {
	case strings.Contains(event, " ") && strings.Contains(event, ","):
		fmt.Println("CASE 3")
	case strings.Contains(event, " "):
		fmt.Println("CASE 1")
	default:
		fmt.Println("UNKNOWN units")
	}
}

func stringifyTagset(m map[string]string) string {
	S := make([]string, 0)
	var x string
	for k, v := range m {
		x = strings.Join([]string{k, "=", v}, "")
		S = append(S, x)
	}
	return strings.Join(S, ",")
}

func TextToTime(timeString string) (time.Time, string) {
	var t time.Time
	var UnixString string
	err := t.UnmarshalText([]byte(timeString))
	if err != nil {
		fmt.Println(err)
		UnixString = "ERROR unmarshalling time object"
	} else {
		UnixString = strconv.Itoa(int(t.UnixNano()))
	}
	return t, UnixString
}

// Use map to generate influx line protocol string
func marshalInfluxLP(m map[string]interface{}, tagset map[string]string, measurement string) {
	// "particles,
	// event=temperature,experiment=CISBAT,location=Archlab,label=thatOne,sample_rate=1000,unit=c
	// value=10e9 0000000000000000000"
	tagsetString := stringifyTagset(tagset)
	_, UnixString := TextToTime(m["published_at"].(string))
	// fields := map[string]string{value: m["data"].(string)}
	s := strings.Join([]string{
		measurement, ",",
		tagsetString, " ",
		m["data"].(string), " ",
		UnixString,
	}, "")
	fmt.Println("\n" + s)
}

func eventCase1() {
	// pass
}

/* MAIN FUNCTION */
func main() {
	// Where the config file is
	credPath := "/Users/eat_sleep_live_skateboarding/Code/go/particle-sse/credentials.txt"

	// The SSE url
	sseURL := "https://api.particle.io/v1/devices/events?access_token="

	// parse values from config
	_map, _ := readLines(credPath)
	settings := parseCreds(_map)

	// Set credentials
	username, ok := settings["user"]
	if ok == false {
		panic("no USER set in config file!")
	}
	password, ok := settings["password"]
	if ok == false {
		panic("no PASSWORD set in config file!")
	}
	database, ok := settings["database"]
	if ok == false {
		panic("no DATABASE set in config file!")
	}
	measurement, ok := settings["measurement"]
	if ok == false {
		panic("no MEASUREMENT set in config file!")
	}
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}

	var input string
	// check devices
	devicesResp := allParticlesCurl(settings["api-key"])
	allParticles := make([]map[string]interface{}, 0)
	json.Unmarshal(devicesResp, &allParticles)

	_particles := AllParticlesMeta{numParticles: len(allParticles)}
	fmt.Printf("The response held: %v values.", _particles.numParticles)

	fmt.Scanln(&input)
	for _, v := range allParticles {
		if v["connected"] == true {
			_particles.numConnected++
		}

	}

	fmt.Printf("\nThere are %v particles connected...\n ", _particles.numConnected)
	/* Pause */
	fmt.Scanln(&input)

	for _, v := range allParticles {
		if v["connected"] == true {
			fmt.Printf("%s %s\n", v["name"], v["id"])
		}
	}

	fmt.Println("\nBeginning SSE Client")
	/* Pause */
	fmt.Scanln(&input)

	// SSE begins here
	sseURL = sseURL + settings["api-key"]
	fmt.Println(sseURL)

	// Create the channel to store SSEresponses
	SSEresp := make(chan string, 2)
	counter := 0
	SSEchanIsReady := make(chan bool)

	// SSEres := make(map[string]interface{})
	sseClient := sse.NewClient(sseURL)
	go sseClient.Subscribe("messages", func(msg *sse.Event) {
		if msg.Event != nil {
			SSEresp <- string(msg.Event)
			SSEchanIsReady <- false
			counter = 0
		} else if msg.Data != nil {
			SSEresp <- string(msg.Data)
			counter = 1
			SSEchanIsReady <- true
		}
	})

	//
	go func() {
		for {
			if <-SSEchanIsReady {
				fmt.Println("\n***")
				res := combineEventAndData(<-SSEresp, <-SSEresp)
				tagset := map[string]string{
					"id":    res["coreid"].(string),
					"event": res["event"].(string),
				}
				fields := map[string]interface{}{value: res["data"].(string)}
				publishedAt, publishedAtString := TextToTime(res["published_at"].(string))
				fmt.Println(publishedAtString)
				marshalInfluxLP(res, tagset, measurement)

				// Create a new point batch
				bp, err := client.NewBatchPoints(client.BatchPointsConfig{
					Database:  database,
					Precision: "s",
				})
				if err != nil {
					log.Fatal(err)
				}

				// this will eventually loop if there were encoded points
				// Create a point and add to batch
				pt, err := client.NewPoint(measurement, tagset, fields, publishedAt)
				if err != nil {
					log.Fatal(err)
				}
				bp.AddPoint(pt)

				// Write the batch
				if err := c.Write(bp); err != nil {
					log.Fatal(err)
				}
				// iterMap(res)
			}
		}
	}()
	// wait for input
	fmt.Scanln(&input)
}

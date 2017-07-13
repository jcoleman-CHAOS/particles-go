package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Decodes an event into: phenom[, units][, numTimes]
func eventSplit(s string) (string, string, int) {
	// Separate the event from the number of times it appears in event.Data
	splitRes := strings.Split(s, ",")
	var numTimes int
	if len(splitRes) > 1 && splitRes[1] != "" {
		// you have numTimesString
		numTimesString := splitRes[1]
		_, err := strconv.Atoi(numTimesString)
		if err != nil {
			fmt.Println(err.Error())
			numTimes = 1
		} else {
			numTimes, _ = strconv.Atoi(numTimesString)
		}
	} else {
		numTimes = 1
	}

	var units string
	// if your event contains units...
	if len(strings.Split(splitRes[0], ".")) > 1 {
		units = strings.Split(splitRes[0], ".")[1]
	} else {
		// if not leave it blank
		units = ""
	}

	// the phenom is what's left after seperation from units
	phenom := strings.Split(splitRes[0], ".")[0]

	return phenom, units, numTimes
}

func broadcastEvent(phenom string, units string, numTimes int) {
	for i := 0; i < numTimes; i++ {
		fmt.Printf("%s %s\n", phenom, units)
	}
}

func decodeEvents(eventString string) ([]string, []string) {
	eventChunks := strings.Split(eventString, " ")
	P := make([]string, 0)
	U := make([]string, 0)
	for i := range eventChunks {
		phenom, unit, numTimes := eventSplit(eventChunks[i])
		for j := 0; j < numTimes; j++ {
			P = append(P, phenom)
			U = append(U, unit)
		}
	}
	return P, U
}

func decodeData(dataString string) {
	dataSplit(dataString)
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return true
	} else {
		return false
	}
}

type rawData struct {
	label  string
	str    string
	hasVal bool
	value  float64
}

func (d *rawData) attributes() {
	if d.hasVal {
		fmt.Printf("Label: %s Str: %s Value: %v\n", d.label, d.str, d.value)
	} else {
		fmt.Printf("Label: %s Str: %s Value: %s\n", d.label, d.str, "nil")
	}
}

func dataSplit(dataString string) []rawData {
	// Split into list that is delimited by spaces
	dataChunks := strings.Split(dataString, " ")
	rawDataPoints := make([]rawData, 0)

	// Iterate through the chunks of strings
	for i := range dataChunks {
		x := strings.Split(dataChunks[i], ":")
		var data rawData
		if len(x) > 1 { // if there is a value attached to the label string
			data.label = x[0]
			var num float64
			if IsNumeric(x[1]) { // if the item attached to the label is a number
				num, _ = strconv.ParseFloat(x[1], 64)
				data.hasVal = true
				data.value = num
			} else { // otherwise it's just a string, but it may still contain a value
				// by convention we denote a string attached to a value with a ";"
				z := strings.Split(x[1], ";")
				if len(z) > 1 { // there are two compnents in the string
					if IsNumeric(z[0]) || IsNumeric(z[1]) { // we have at least one val
						data.hasVal = true
						// check if we have one numeric one string
						if IsNumeric(z[0]) { // the first val is a number...
							// we assume the second to be a string
							data.value, _ = strconv.ParseFloat(z[0], 64)
							data.str = z[1]
						} else { // the second val is a number...
							// we assume the first to be a string
							data.str = z[0]
							data.value, _ = strconv.ParseFloat(z[1], 64)
						}
					} else { // we just have two strings
						data.str = strings.Join([]string{z[0], z[1]}, " ")
					}
				} else { // there is a str, but no numeric
					data.str = x[1]
				}
			}
		} else { // otherwise presume you just have a value, no label
			if IsNumeric(x[0]) {
				num, _ := strconv.ParseFloat(x[0], 64)
				data.hasVal = true
				data.value = num
			} else {
				data.str = x[0]
			}
		}
		rawDataPoints = append(rawDataPoints, data)
	}
	return rawDataPoints
}

func matchEventsAndData(phenoms []string, units []string, points []rawData) {
	// Several cases:
	// 1. phenoms, units, and points are all same length
	if len(phenoms) == len(units) && len(phenoms) == len(points) {
		for i := range phenoms {
			d := points[i]
			fmt.Printf("Label: %s Phenom: %s Unit: %s Value: %v String: %s \n",
				d.label, phenoms[i], units[i], d.value, d.str)
		}
	} else { // For now print an error message
		fmt.Println("there was an error in the input string:")
		fmt.Println("\tNumber of Phenoms, Units and Values do not match.")
	}
}

func main() {
	eventString := "temp"
	dataString := "10"
	P, U := decodeEvents(eventString)
	// fmt.Println(strings.Join(P, ","))
	// fmt.Println(strings.Join(U, ","))
	D := dataSplit(dataString)
	// fmt.Printf("There are %v PHENOMS.", len(P))
	// fmt.Printf("There are %v UNITS.", len(U))
	// fmt.Printf("There are %v DATAPOINTS.\n", len(D))
	matchEventsAndData(P, U, D)
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type StationModel struct {
	Type             string `json:"type"`
	StationName      string `json:"stationName"`
	StationShortCode string `json:"stationShortCode"`
}

// global list of stations to prevent refetch
var stations []StationModel

func getStations() []StationModel {
	// use cached version if available
	if stations != nil {
		return stations
	}

	// fetch station info from API
	url := fmt.Sprintf("%s/%s", apiUrl, "metadata/stations")
	res, getError := http.Get(url)

	// if fetching fails...
	if getError != nil {
		fmt.Println("Could not fetch station list.")
		os.Exit(1)
	}

	// read all the data sent to byte array and handle error case
	body, readError := io.ReadAll(res.Body)
	if readError != nil {
		fmt.Println("Could not read station list.")
		os.Exit(1)
	}

	// write to file
	writeJsonToFile("station.json", string(body))

	// as data is returned as an array with one entry, Unmarshal and return the first entry
	jsonError := json.Unmarshal(body, &stations)
	if jsonError != nil {
		fmt.Println("Could not parse station list.")
		os.Exit(1)
	}

	// store to cache
	return stations
}

func getStationByShortCode(stationShortCode string) StationModel {
	stations := getStations()
	for _, station := range stations {
		if station.StationShortCode == stationShortCode {
			return station
		}
	}
	return stations[0]
}

func (s StationModel) print() {
	fmt.Printf("%s: %s\n", s.StationShortCode, s.StationName)
}

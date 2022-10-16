package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const apiUrl = "https://rata.digitraffic.fi/api/v1"

type TrainModel struct {
	TrainNumber   int                 `json:"trainNumber"`
	TrainType     string              `json:"trainType"`
	Version       int                 `json:"version"`
	TimeTableRows []TimeTableRowModel `json:"timeTableRows"`
}

func getTrain(trainNumber int) TrainModel {

	// fetch train info from API
	url := fmt.Sprintf("%s/%s/%s/%s", apiUrl, "trains", getTimestamp(), strconv.Itoa(trainNumber))
	res, getError := http.Get(url)

	// if fetching fails...
	if getError != nil {
		fmt.Println("Could not fetch train info.")
		os.Exit(1)
	}

	// read all the data sent to byte array and handle error case
	body, readError := io.ReadAll(res.Body)
	if readError != nil {
		fmt.Println("Could not read train info.")
		os.Exit(1)
	}

	// as data is returned as an array with one entry, Unmarshal and return the first entry
	var trains []TrainModel
	jsonError := json.Unmarshal(body, &trains)
	if jsonError != nil {
		fmt.Println("Could not parse train info.")
		os.Exit(1)
	}
	return trains[0]
}

func getAllTrains() []TrainModel {
	// fetch train info from API
	url := fmt.Sprintf("%s/%s/%s", apiUrl, "trains", getTimestamp())
	fmt.Println(url)
	res, getError := http.Get(url)

	// if fetching fails...
	if getError != nil {
		fmt.Println("Could not fetch trains list.")
		os.Exit(1)
	}

	// read all the data sent to byte array and handle error case
	body, readError := io.ReadAll(res.Body)
	if readError != nil {
		fmt.Println("Could not read trains list.")
		os.Exit(1)
	}

	// as data is returned as an array with one entry, Unmarshal and return the first entry
	var trains []TrainModel
	jsonError := json.Unmarshal(body, &trains)
	if jsonError != nil {
		fmt.Println("Could not parse trains list.")
		os.Exit(1)
	}
	return trains
}

func (t TrainModel) printName() {
	fmt.Println("Train", t.TrainNumber, "("+t.TrainType+")", "["+strconv.Itoa(t.Version)+"]")
}

func (t TrainModel) printTimeTableRows() {
	for _, ttr := range t.TimeTableRows {
		ttr.print()
	}
}

func (t TrainModel) print() {
	t.printName()
	t.printTimeTableRows()
}

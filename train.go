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
	TrainNumber int    `json:"trainNumber"`
	TrainType   string `json:"trainType"`
}

func getTrain(trainNumber int) TrainModel {
	// fetch train info from API
	url := apiUrl + "/trains/latest/" + strconv.Itoa(trainNumber)
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

func (t TrainModel) print() {
	fmt.Println("Train", t.TrainNumber, "is type", t.TrainType)
}

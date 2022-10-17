package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
)

const apiUrl = "https://rata.digitraffic.fi/api/v1"

type TrainModel struct {
	TrainNumber    int                 `json:"trainNumber"`
	TrainType      string              `json:"trainType"`
	CommuterLineID string              `json:"commuterLineID"`
	Version        int                 `json:"version"`
	TimeTableRows  []TimeTableRowModel `json:"timeTableRows"`
	Cancelled      bool                `json:"cancelled"`
}

// a generic wrapper implementing request to train api
func trainApi(path string, logfile string) []TrainModel {
	url := fmt.Sprintf("%s/%s", apiUrl, path)
	res, getError := http.Get(url)

	// if fetching fails...
	if getError != nil {
		fmt.Println("Request " + path + " failed.")
		os.Exit(1)
	}

	// read all the data sent to byte array and handle error case
	body, readError := io.ReadAll(res.Body)
	if readError != nil {
		fmt.Println("Failed to receive proper data from " + path)
		os.Exit(1)
	}

	// write to file
	writeJsonToFile(logfile, string(body))

	// unmarshal and return the whole array
	var trains []TrainModel
	jsonError := json.Unmarshal(body, &trains)
	if jsonError != nil {
		fmt.Println("Could not parse " + path + " response")
		os.Exit(1)
	}
	return trains
}

// fetch a single train info
func getTrain(trainNumber int) TrainModel {
	trainNumberAsString := strconv.Itoa(trainNumber)
	path := fmt.Sprintf("%s/%s/%s", "trains", getTimestamp(), trainNumberAsString)
	logfile := fmt.Sprintf("%s.json", trainNumberAsString)
	trains := trainApi(path, logfile)
	return trains[0]
}

// fetch list of all trains (today)
func getAllTrains() []TrainModel {
	path := fmt.Sprintf("%s/%s", "trains", getTimestamp())
	trains := trainApi(path, "all.json")
	return trains
}

// get all trains stopping station (today, not yet stopped)
func getTrainsByStation(stationShortCode string) []TrainModel {
	path := fmt.Sprintf("live-trains/station/%s?include_nonstopping=false&departing_trains=10", stationShortCode)
	logfile := fmt.Sprintf("%s.json", stationShortCode)
	trains := trainApi(path, logfile)

	// sort by ScheduledTime
	sort.SliceStable(trains, func(a, b int) bool {
		trainA := trains[a]
		trainB := trains[b]

		entryA := trainA.getStationEntry(stationShortCode)
		entryB := trainB.getStationEntry(stationShortCode)

		timeA := parseToTime(entryA.ScheduledTime)
		timeB := parseToTime(entryB.ScheduledTime)

		return timeA.Before(timeB)
	})
	return trains
}

func (t TrainModel) printTimeTableRows() {
	for index, ttr := range t.TimeTableRows {
		if ttr.isStopping() {
			// print departing rows, except the final destination will be printed as well
			isFinalDestionation := ttr.isArrival() && index == len(t.TimeTableRows)-1
			if ttr.isDeparture() || isFinalDestionation {
				ttr.print()
			}
		}
	}
}

func (t TrainModel) isCancelled() bool {
	return t.Cancelled
}

func (t TrainModel) printHeader() {
	fmt.Printf("%s\n", t.getHeader())
}

func (t TrainModel) getStationEntry(stationShortCode string) TimeTableRowModel {
	// find departure record for this station
	for _, row := range t.TimeTableRows {
		if row.StationShortCode == stationShortCode && row.isDeparture() {
			return row
		}
	}
	return t.TimeTableRows[0]
}

func (t TrainModel) printScheduleEntry(stationShortCode string) {
	// get station entry
	ttr := t.getStationEntry(stationShortCode)

	// get scheduled time as string
	departureTime := parseToTime(ttr.ScheduledTime).Format("15:04")

	// define some wording based on is the train yet stopped
	estimatedDeparture := yellow("estimated departure at")
	if ttr.hasStopped() {
		estimatedDeparture = green("departured at")
	}

	// print it all
	fmt.Printf("%s to %s %s %s\n", t.getHeader(), t.getFinalDestination(), estimatedDeparture, departureTime)
}

func (t TrainModel) getHeader() string {
	return fmt.Sprintf("%s Train %s", t.getType(), t.getNumber())
}
func (t TrainModel) getNumber() string {
	return strconv.Itoa(t.TrainNumber)
}

func (t TrainModel) getType() string {
	if t.CommuterLineID != "" {
		return t.CommuterLineID
	}
	return t.TrainType
}

func (t TrainModel) getFinalDestination() string {
	lastDestination := t.TimeTableRows[len(t.TimeTableRows)-1]
	return lastDestination.getStationName()
}

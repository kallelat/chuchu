package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
)

const apiUrl = "https://rata.digitraffic.fi/api/v1"

type TrainModel struct {
	TrainNumber      int                 `json:"trainNumber"`
	TrainType        string              `json:"trainType"`
	CommuterLineID   string              `json:"commuterLineID"`
	Version          int                 `json:"version"`
	TimeTableRows    []TimeTableRowModel `json:"timeTableRows"`
	Cancelled        bool                `json:"cancelled"`
	RunningCurrently bool                `json:"runningCurrently"`
}

// a generic wrapper implementing request to train api
func trainApi(path string, logfile string) ([]TrainModel, error) {
	url := fmt.Sprintf("%s/%s", apiUrl, path)
	res, getError := http.Get(url)

	// if fetching fails...
	if getError != nil {
		return nil, getError
	}

	// read all the data sent to byte array and handle error case
	body, readError := io.ReadAll(res.Body)
	if readError != nil {
		return nil, readError
	}

	// write to file
	writeJsonToFile(logfile, string(body))

	// unmarshal and return the whole array
	var trains []TrainModel
	jsonError := json.Unmarshal(body, &trains)
	if jsonError != nil {
		return nil, jsonError
	}
	return trains, nil
}

// fetch a single train info
func getTrain(trainNumber int) (*TrainModel, error) {
	trainNumberAsString := strconv.Itoa(trainNumber)
	path := fmt.Sprintf("%s/%s/%s", "trains", getTimestamp(), trainNumberAsString)
	logfile := fmt.Sprintf("%s.json", trainNumberAsString)
	trains, err := trainApi(path, logfile)
	if err != nil {
		return nil, err
	}

	if len(trains) == 0 {
		return nil, errors.New("NoTrainsFound")
	}
	return &trains[0], nil
}

// fetch list of all trains (today)
func getAllTrains() ([]TrainModel, error) {
	path := fmt.Sprintf("%s/%s", "trains", getTimestamp())
	trains, err := trainApi(path, "all.json")
	if err != nil {
		return nil, err
	}
	return trains, nil
}

// get all trains stopping station (today, not yet stopped)
func getTrainsByStation(stationShortCode string) ([]TrainModel, error) {
	path := fmt.Sprintf("live-trains/station/%s?include_nonstopping=false&departing_trains=10", stationShortCode)
	logfile := fmt.Sprintf("%s.json", stationShortCode)
	trains, err := trainApi(path, logfile)

	if err != nil {
		return nil, err
	}

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
	return trains, nil
}

func (t TrainModel) isCancelled() bool {
	return t.Cancelled
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

func (t TrainModel) toScheduleEntry(stationShortCode string) string {
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
	return fmt.Sprintf("%s to %s %s %s", t.toString(), t.getFinalDestination(), estimatedDeparture, departureTime)
}

func (t TrainModel) toString() string {
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

func (t TrainModel) hasDeparted() bool {
	// get the first timetable row (it should only have DEPARTURE item)
	ttr := t.TimeTableRows[0]

	// if the train has actula departing time added, then the train has left
	return ttr.ActualTime != ""
}

func (t TrainModel) getDelayInMinutes() int {
	var lastTtr TimeTableRowModel

	// check every timetable row
	for index, ttr := range t.TimeTableRows {
		if ttr.isStopping() {
			// when finding the first station that the train hasn't stopped yet, check delay from previous station
			if !ttr.hasStopped() && index != 0 {
				return lastTtr.DifferenceInMinutes
			}

			// save reference
			lastTtr = ttr
		}
	}

	// default to zero
	return 0
}

func (t TrainModel) getNextStop() TimeTableRowModel {
	// check every timetable row
	for _, ttr := range t.TimeTableRows {
		if ttr.isStopping() {
			// when finding the first station that the train hasn't stopped yet, check delay from previous station
			if !ttr.hasStopped() {
				return ttr
			}
		}
	}

	// default to first actual stop
	return t.TimeTableRows[1]
}

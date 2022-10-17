package main

import (
	"fmt"
	"strconv"
)

type TimeTableRowModel struct {
	Type                string `json:"type"`
	StationShortCode    string `json:"stationShortCode"`
	Cancelled           bool   `json:"cancelled"`
	ScheduledTime       string `json:"scheduledTime"`
	ActualTime          string `json:"actualTime"`
	DifferenceInMinutes int    `json:"differenceInMinutes"`
	TrainStopping       bool   `json:"trainStopping"`
}

func (ttr TimeTableRowModel) print() {
	// get scheduled time as string
	scheduledTime := parseToTime(ttr.ScheduledTime).Format("15:04")

	// define some wording based on is the train yet stopped
	arrivedOrEstimated := yellow("estimated arrival at")
	if ttr.hasStopped() {
		arrivedOrEstimated = green("arrived at")
	}

	// main output
	output := fmt.Sprintf(" => %s [%s], %s %s", ttr.getStationName(), ttr.StationShortCode, arrivedOrEstimated, scheduledTime)

	// if delayed, add notification
	if ttr.DifferenceInMinutes > 0 {
		output += red(fmt.Sprintf(" (delayed +%s minutes)", strconv.Itoa(ttr.DifferenceInMinutes)))
	}

	// print
	fmt.Println(output)
}

func (ttr TimeTableRowModel) getStationName() string {
	station := getStationByShortCode(ttr.StationShortCode)
	return station.StationName
}

func (ttr TimeTableRowModel) isStopping() bool {
	return ttr.TrainStopping
}

func (ttr TimeTableRowModel) isArrival() bool {
	return ttr.Type == "ARRIVAL"
}

func (ttr TimeTableRowModel) isDeparture() bool {
	return ttr.Type == "DEPARTURE"
}

func (ttr TimeTableRowModel) hasStopped() bool {
	return ttr.ActualTime != ""
}

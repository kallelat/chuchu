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
	scheduledTime := parseToTime(ttr.ScheduledTime).Format("15:04")
	output := fmt.Sprintf(" ==> %s, arrival at %s", ttr.getStationName(), scheduledTime)
	if ttr.DifferenceInMinutes > 0 {
		output += red(fmt.Sprintf(" delayed %s minutes", strconv.Itoa(ttr.DifferenceInMinutes)))
	}
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

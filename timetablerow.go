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
}

func (ttr TimeTableRowModel) print() {
	output := fmt.Sprintf(" ==> %s, arrival at %s, delay (%s minutes delayed)", ttr.getStationName(), ttr.ScheduledTime, strconv.Itoa(ttr.DifferenceInMinutes))
	fmt.Println(output)
}

func (ttr TimeTableRowModel) getStationName() string {
	station := getStationByShortCode(ttr.StationShortCode)
	return station.StationName
}

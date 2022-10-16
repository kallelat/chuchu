package main

import "fmt"

type TimeTableRowModel struct {
	Type                string `json:"type"`
	StationShortCode    string `json:"stationShortCode"`
	Cancelled           bool   `json:"cancelled"`
	ScheduledTime       string `json:"scheduledTime"`
	ActualTime          string `json:"actualTime"`
	DifferenceInMinutes int    `json:"differenceInMinutes"`
}

func (ttr TimeTableRowModel) print() {
	fmt.Println("  @"+ttr.StationShortCode+":", ttr.DifferenceInMinutes, "delayed")
}

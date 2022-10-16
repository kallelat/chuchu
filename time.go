package main

import (
	"fmt"
	"os"
	"time"
)

const (
	HHMM     = "hh:mm"
	YYYYMMDD = "2006-01-02"
)

func getTimestamp() string {
	now := time.Now().UTC()
	return now.Format(YYYYMMDD)
}

func parseToTime(value string) time.Time {
	time, err := time.Parse(time.RFC3339, value)
	if err != nil {
		fmt.Println("Can't parse date")
		os.Exit(1)
	}
	return time
}

package main

import "time"

const (
	YYYYMMDD = "2006-01-02"
)

func getTimestamp() string {
	now := time.Now().UTC()
	return now.Format(YYYYMMDD)
}

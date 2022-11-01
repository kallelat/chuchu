package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
)

const port = 8080

func serve() {
	// register handlers
	http.HandleFunc("/station/", func(w http.ResponseWriter, r *http.Request) {
		stationCode := path.Base(r.RequestURI)
		stationResponse(w, stationCode)
	})

	http.HandleFunc("/train/", func(w http.ResponseWriter, r *http.Request) {
		trainNumber := path.Base(r.RequestURI)
		trainResponse(w, trainNumber)
	})

	portAsString := strconv.Itoa(port)
	fmt.Printf("Starting server at port %s\n", portAsString)
	if err := http.ListenAndServe(":"+portAsString, nil); err != nil {
		log.Fatal(err)
	}
}

func stationResponse(w http.ResponseWriter, stationCode string) {
	setHTMLType(w)

	// print station name
	writeLine(w, "Station: "+stationCode)

	// get trains
	trains, err := getTrainsByStation(stationCode)
	if err != nil {
		errorResponse(w, "Can't parse station: "+stationCode)
		return
	}

	// print link for each train
	for _, train := range trains {
		ttr := train.getStationEntry(stationCode)
		time := parseToTime(ttr.ScheduledTime).Format("15:04")
		str := fmt.Sprintf("  <a href=\"/train/%s\">%s</a> [%s]<br/>", train.getNumber(), train.toString(), time)
		fmt.Fprint(w, str)
	}
}

func trainResponse(w http.ResponseWriter, trainNumber string) {
	setHTMLType(w)

	// parse as int
	trainNumberAsInt, err := strconv.Atoi(trainNumber)

	// throw error page
	if err != nil {
		errorResponse(w, "Can't parse train number: "+trainNumber)
		return
	}

	// fetch train info
	train, err := getTrain(trainNumberAsInt)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}

	// print train header
	writeLine(w, train.toString())

	// print rows
	for _, ttr := range train.TimeTableRows {
		if ttr.isDeparture() {
			writeLine(w, ttr.toHTMLString())
		}
	}
}

func errorResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(501)
	fmt.Fprintln(w, message)
}

func writeLine(w http.ResponseWriter, str string) {
	finalStr := fmt.Sprintf("%s<br/>", str)
	fmt.Fprintln(w, finalStr)
}
func setHTMLType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

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
		station := path.Base(r.RequestURI)
		stationResponse(w, station)
	})

	http.HandleFunc("/train/", func(w http.ResponseWriter, r *http.Request) {
		train := path.Base(r.RequestURI)
		trainResponse(w, train)
	})

	portAsString := strconv.Itoa(port)
	fmt.Printf("Starting server at port %s\n", portAsString)
	if err := http.ListenAndServe(":"+portAsString, nil); err != nil {
		log.Fatal(err)
	}
}

func stationResponse(w http.ResponseWriter, stationName string) {
	fmt.Fprint(w, "Station: "+stationName)
}

func trainResponse(w http.ResponseWriter, trainNumber string) {
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
	fmt.Fprintln(w, train.getHeader())
}

func errorResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(501)
	fmt.Fprintln(w, message)
}

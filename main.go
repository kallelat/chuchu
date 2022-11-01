package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// pick the args, later provide support for other attributes as well
	trainNumberAttribute := flag.Int("train", 0, "train number as integer, usage: -train <trainNumber>")
	allTrainsAttribute := flag.Bool("all", false, "lists all trains currently available, usage: -all")
	watchTrainsAttribute := flag.Int("watch", 0, "watch a certain train and let user know if there are changes, usage: -watch <trainNumber>")
	stationsAttribute := flag.Bool("stations", false, "list all stations")
	stationAttribute := flag.String("station", "", "list trains today by station, usage -station <stationShortCode>")
	serverAttribute := flag.Bool("server", false, "starts a server user can use to poll train schedules")

	flag.Parse()

	if *serverAttribute {
		serve()
	} else if *watchTrainsAttribute != 0 {
		watch(*watchTrainsAttribute)
		os.Exit(1)
	} else if *trainNumberAttribute != 0 {
		train, err := getTrain(*trainNumberAttribute)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// print train header
		fmt.Println(train.toString())

		// if cancelled, print status and exit
		if train.isCancelled() {
			fmt.Println("TRAIN IS CANCELLED!")
			os.Exit(0)
		}

		// if not cancelled, print train timetablerows
		for index, ttr := range train.TimeTableRows {
			if ttr.isStopping() {
				// print departing rows, except the final destination will be printed as well
				isFinalDestionation := ttr.isArrival() && index == len(train.TimeTableRows)-1
				if ttr.isDeparture() || isFinalDestionation {
					fmt.Println(ttr.toString())
				}
			}
		}
	} else if *allTrainsAttribute {
		trains, err := getAllTrains()

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// print all train headers
		for _, train := range trains {
			fmt.Println(train.toString())
		}
	} else if *stationAttribute != "" {
		trains, err := getTrainsByStation(*stationAttribute)

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// print schedule for each train stopping in the stations
		for _, train := range trains {
			str := train.toScheduleEntry(*stationAttribute)
			fmt.Println(str)
		}
	} else if *stationsAttribute {
		stations := getStations()

		for _, station := range stations {
			fmt.Println(station.toString())
		}
	}
}

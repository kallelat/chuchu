package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func watch(trainNumber int) {
	// start by fetching the train info
	train, err := getTrain(trainNumber)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("Watching: %s\n", train.getHeader())
	for {
		// first, cancellation check
		if train.isCancelled() {
			text := "cancelled!"
			fmt.Printf(" ==> %s\n", red(text))
			break // exit loop and quite app
		}

		// check if the train has not even left yet
		if train.hasDeparted() {
			// get the current delay
			delay := train.getDelayInMinutes()
			nextStop := train.getNextStop()
			stationName := nextStop.getStationName()
			scheduledTime := parseToTime(nextStop.ScheduledTime).Format("15:04")
			if delay == 0 {
				text := "on schedule!"
				fmt.Printf(" ==> %s [next stop %s at %s]\n", green(text), stationName, scheduledTime)
			} else {
				text := fmt.Sprintf(" ==> delayed by %s minutes! [next stop %s at %s]", strconv.Itoa(delay), stationName, scheduledTime)
				fmt.Printf(" ==> %s\n", red(text))
			}
		} else {
			text := "hasn't departed yet!"
			fmt.Printf(" ==> %s", yellow(text))
		}

		// update train info and sleep before next iteration
		train, err = getTrain(trainNumber)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// pause execution for 15 sec
		time.Sleep(15 * time.Second)
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// pick the args, later provide support for other attributes as well
	trainNumberAttribute := flag.Int("train", 0, "train number as integer")
	allTrainsAttribute := flag.Bool("all", false, "lists all trains currently available")

	flag.Parse()

	if *trainNumberAttribute != 0 {
		train := getTrain(*trainNumberAttribute)
		train.printName()

		// if cancelled, print status and exit
		if train.isCancelled() {
			fmt.Println("TRAIN IS CANCELLED!")
			os.Exit(0)
		}

		// if not cancelled, print train timetablerows
		train.printTimeTableRows()
	}

	if *allTrainsAttribute {
		trains := getAllTrains()
		for _, train := range trains {
			train.printName()
		}

	}
}

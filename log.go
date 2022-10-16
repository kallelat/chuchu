package main

import (
	"fmt"
	"os"
)

func writeJsonToFile(filename string, json string) {
	// create a file, only truncates if already exists
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Can't write json outout")
		os.Exit(1)
	}

	// write json to the file
	defer f.Close()
	_, err2 := f.WriteString(json)
	if err2 != nil {
		fmt.Println("Writing json output failed")
	}
}

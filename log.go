package main

import (
	"fmt"
	"os"
)

func writeJsonToFile(filename string, json string) {
	f, err := os.Create(filename)

	if err != nil {
		fmt.Println("Can't write json outout")
		os.Exit(1)
	}

	defer f.Close()

	_, err2 := f.WriteString(json)

	if err2 != nil {
		fmt.Println("Writing json output failed")
	}
}

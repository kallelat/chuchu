package main

import "fmt"

const (
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
)

// a generic function to change text color, resets color at the end
func color(color string, text string) string {
	return fmt.Sprintf("%s%s\033[0m", color, text)
}

func red(text string) string {
	return color(RED, text)
}

func green(text string) string {
	return color(GREEN, text)
}

func yellow(text string) string {
	return color(YELLOW, text)
}

package main

import "fmt"

func red(text string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", text)
}

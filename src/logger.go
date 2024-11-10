package main

import "fmt"

func log(ty, format string, args ...any) {
	fmt.Printf(ty+format, args...)
}

func logError(format string, args ...any) {
	log("error: ", format, args...)
}

func logInfo(format string, args ...any) {
	log("info: ", format, args...)
}

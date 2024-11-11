package main

import "fmt"

var CanLog = false

func log(ty, format string, args ...any) {
	fmt.Printf(ty+format, args...)
}

func logError(format string, args ...any) {
	log("error: ", format, args...)
}

func logInfo(format string, args ...any) {
	if CanLog {
		log("info: ", format, args...)
	}
}

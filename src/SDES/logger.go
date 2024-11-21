package sdes

import "fmt"

var CanLog = false

func log(ty, format string, args ...any) {
	fmt.Printf(ty+format, args...)
}

func LogError(format string, args ...any) {
	log("error: ", format, args...)
}

func LogInfo(format string, args ...any) {
	if CanLog {
		log("info: ", format, args...)
	}
}

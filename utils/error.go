package utils

import "log"

func CheckFatal(err error, msgs ...string) {
	if err == nil {
		return
	}

	log.Fatal(err, msgs)
}

func CheckPanic(err error, msgs ...string) {
	if err == nil {
		return
	}

	log.Panic(err, msgs)
}

func CheckWarn(err error, msgs ...string) {
	if err == nil {
		return
	}

	log.Print(err, "WARN ", msgs)
}

func LogInfo(err error, msgs ...string) {
	if err == nil {
		return
	}

	log.Print(err, "INFO ", msgs)
}

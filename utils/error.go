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

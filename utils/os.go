package utils

import (
	"path/filepath"
	"runtime"
)

func GetCurrentFile() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}

func GetCurrrentFolder() string {
	_, filename, _, _ := runtime.Caller(1)
	directory := filepath.Dir(filename)
	return directory
}

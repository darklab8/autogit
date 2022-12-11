package utils

import (
	"path/filepath"
	"runtime"
)

func GetCurrrentFolder(folder_name string) string {
	_, filename, _, _ := runtime.Caller(1)
	directory := filepath.Dir(filename)
	return directory
}

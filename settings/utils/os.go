package utils

import (
	"autogit/settings/types"
	"path/filepath"
	"runtime"
)

func GetCurrentFile() types.FilePath {
	_, filename, _, _ := runtime.Caller(1)
	return types.FilePath(filename)
}

func GetCurrrentFolder() types.FilePath {
	_, filename, _, _ := runtime.Caller(1)
	directory := filepath.Dir(filename)
	return types.FilePath(directory)
}

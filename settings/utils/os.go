package utils

import (
	"autogit/settings/logus"
	"autogit/settings/types"
	"os"
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

func GetProjectDir() types.FilePath {
	path, err := os.Getwd()
	if folder_override, ok := os.LookupEnv("AUTOGIT_PROJECT_FOLDER"); ok {
		path = folder_override
	}
	logus.CheckFatal(err, "unable to get workdir")
	return types.FilePath(path)
}

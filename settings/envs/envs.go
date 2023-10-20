package envs

import (
	"autogit/settings/types"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/*
during unit tests, code grabs wrong folder because
unit tests are located in nested folders.
And autogit is able to run with correct settings only if run from project root
TODO fix actually to detect root folder of it, then it will not be necessary value
*/
var TestProjectFolder types.ProjectFolder

var LogTurnJSONLogging bool
var LogShowFileLocations bool
var LogLevel types.LogLevel

var PathUserHome types.FilePath
var PathGitConfig types.FilePath

const (
	EnvTrue = "true"
)

func init() {
	TestProjectFolder = types.ProjectFolder(os.Getenv("AUTOGIT_TEST_PROJECT_FOLDER"))
	LogTurnJSONLogging = strings.ToLower(os.Getenv("AUTOGIT_LOG_JSON")) == EnvTrue
	LogShowFileLocations = strings.ToLower(os.Getenv("AUTOGIT_LOG_SHOW_FILE_LOCATIONS")) == EnvTrue

	log_level_str, is_log_level_set := os.LookupEnv("AUTOGIT_LOG_LEVEL")
	if !is_log_level_set {
		log_level_str = "INFO"
	}
	LogLevel = types.LogLevel(log_level_str)

	PathUserHome = types.FilePath(os.Getenv("HOME"))

	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Unable to get UserHomeDir, e=%v", err)
	}
	PathUserHome = types.FilePath(dirname)
	PathGitConfig = types.FilePath(filepath.Join(string(PathUserHome), ".gitconfig"))
}

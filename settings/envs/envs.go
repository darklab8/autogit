package envs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/darklab8/darklab_goutils/goutils/logus_core/logus_types"
	"github.com/darklab8/darklab_goutils/goutils/utils/utils_types"
)

/*
during unit tests, code grabs wrong folder because
unit tests are located in nested folders.
And autogit is able to run with correct settings only if run from project root
TODO fix actually to detect root folder of it, then it will not be necessary value
*/

var LogTurnJSONLogging bool
var LogShowFileLocations bool
var LogLevel logus_types.LogLevel

var PathUserHome utils_types.FilePath
var PathGitConfig utils_types.FilePath

const (
	EnvTrue = "true"
)

func init() {
	LogTurnJSONLogging = strings.ToLower(os.Getenv("AUTOGIT_LOG_JSON")) == EnvTrue
	LogShowFileLocations = strings.ToLower(os.Getenv("AUTOGIT_LOG_SHOW_FILE_LOCATIONS")) == EnvTrue

	log_level_str, is_log_level_set := os.LookupEnv("AUTOGIT_LOG_LEVEL")
	if !is_log_level_set {
		log_level_str = "INFO"
	}
	LogLevel = logus_types.LogLevel(log_level_str)

	PathUserHome = utils_types.FilePath(os.Getenv("HOME"))

	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Unable to get UserHomeDir, e=%v", err)
	}
	PathUserHome = utils_types.FilePath(dirname)
	PathGitConfig = utils_types.FilePath(filepath.Join(string(PathUserHome), ".gitconfig"))
}

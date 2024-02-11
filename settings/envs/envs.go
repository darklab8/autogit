package envs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

/*
during unit tests, code grabs wrong folder because
unit tests are located in nested folders.
And autogit is able to run with correct settings only if run from project root
TODO fix actually to detect root folder of it, then it will not be necessary value
*/

var PathUserHome utils_types.FilePath
var PathGitConfig utils_types.FilePath

func init() {
	PathUserHome = utils_types.FilePath(os.Getenv("HOME"))

	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Unable to get UserHomeDir, e=%v", err)
	}
	PathUserHome = utils_types.FilePath(dirname)
	PathGitConfig = utils_types.FilePath(filepath.Join(string(PathUserHome), ".gitconfig"))
}

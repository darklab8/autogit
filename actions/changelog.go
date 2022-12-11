package actions

import (
	"autogit/git"
	"fmt"
)

func Changelog() string {
	return fmt.Sprintf("%v", (&git.Repository{}).GetRepoInWorkDir().GetLogs(git.HEAD_Hash))
}

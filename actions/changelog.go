package actions

import (
	"autogit/git"
	"fmt"
)

func Changelog() string {
	return fmt.Sprintf("%v", (&git.Repository{}).NewRepoInWorkDir().GetLogs(git.HEAD_Hash))
}

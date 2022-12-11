package changelog

import (
	"autogit/git"
)

func Run() {

	(&git.Repository{}).GetRepoInWorkDir().GetLogs(git.HEAD_Hash)
}

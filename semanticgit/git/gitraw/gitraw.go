package gitraw

import (
	"github.com/darklab8/autogit/v2/settings/logus"
	"github.com/darklab8/go-utils/utils/utils_os"

	"github.com/go-git/go-git/v5"
)

func NewGitRepo() *git.Repository {
	path := utils_os.GetProjectDir()

	repo, err := git.PlainOpenWithOptions(string(path), &git.PlainOpenOptions{DetectDotGit: true})
	logus.Log.CheckFatal(err, "unable to open git")
	return repo
}

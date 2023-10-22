package gitraw

import (
	"autogit/settings/logus"
	"autogit/settings/utils"

	"github.com/go-git/go-git/v5"
)

func NewGitRepo() *git.Repository {
	path := utils.GetProjectDir()

	repo, err := git.PlainOpenWithOptions(string(path), &git.PlainOpenOptions{DetectDotGit: true})
	logus.CheckFatal(err, "unable to open git")
	return repo
}

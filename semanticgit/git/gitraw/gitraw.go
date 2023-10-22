package gitraw

import (
	"autogit/settings/logus"
	"os"

	"github.com/go-git/go-git/v5"
)

func NewGitRepo() *git.Repository {
	path, err := os.Getwd()
	if folder_override, ok := os.LookupEnv("AUTOGIT_PROJECT_FOLDER"); ok {
		path = folder_override
	}

	logus.CheckFatal(err, "unable to get workdir")

	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	logus.CheckFatal(err, "unable to open git")
	return repo
}

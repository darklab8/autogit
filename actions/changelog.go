package actions

import (
	"autogit/actions/changelog"
	"autogit/git"
	sGit "autogit/parser/semanticGit"
)

func Changelog() string {
	g := (&sGit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())
	return changelog.ChangelogData{}.New(g).Render()
}

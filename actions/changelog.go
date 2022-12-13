package actions

import (
	"autogit/actions/changelog"
	"autogit/git"
	sGit "autogit/parser/semanticGit"
)

var ChangelogTag *string

func Changelog() string {
	g := (&sGit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())
	return changelog.ChangelogData{Tag: *ChangelogTag}.New(g).Render()
}

package actions

import (
	"autogit/actions/changelog"
	"autogit/semanticgit"
	"autogit/semanticgit/git"
)

var ChangelogTag *string

func Changelog() string {
	g := (&semanticgit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())
	return changelog.ChangelogData{Tag: *ChangelogTag}.New(g).Render()
}

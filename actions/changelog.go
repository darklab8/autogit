package actions

import (
	"autogit/actions/changelog"
	"autogit/actions/validation"
	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"autogit/utils"
)

var ChangelogTag *string
var ChangelogValidate *bool

func Changelog() string {
	g := (&semanticgit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())
	rendered_changelog := changelog.ChangelogData{Tag: *ChangelogTag}.New(g).Render()

	if *ChangelogValidate {
		log_records := g.GetChangelogByTag(*ChangelogTag, false)
		for _, record := range log_records {
			err := validation.Validate(&record)
			utils.CheckFatal(err)
		}
	}

	return rendered_changelog
}

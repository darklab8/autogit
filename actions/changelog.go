package actions

import (
	"autogit/actions/changelog"
	"autogit/actions/validation"
	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"autogit/semanticgit/semver"
	"autogit/utils"
)

var ChangelogTag *string
var ChangelogValidate *bool

var ChangelogDisableVFlag *bool
var ChangelogBuildMeta *string
var ChangelogAlpha *bool
var ChangelogBeta *bool
var ChangelogPrerelease *bool
var ChangelogPublish *bool

func Changelog() string {
	g := (&semanticgit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())
	rendered_changelog := changelog.ChangelogData{Tag: *ChangelogTag}.New(g, semver.OptionsSemVer{
		DisableVFlag:  *ChangelogDisableVFlag,
		EnableNewline: false,
		Build:         *ChangelogBuildMeta,
		Alpha:         *ChangelogAlpha,
		Beta:          *ChangelogBeta,
		Rc:            *ChangelogPrerelease,
		Publish:       *ChangelogPublish,
	}).Render()

	if *ChangelogValidate {
		log_records := g.GetChangelogByTag(*ChangelogTag, false)
		for _, record := range log_records {
			err := validation.Validate(&record)
			utils.CheckFatal(err)
		}
	}

	return rendered_changelog
}

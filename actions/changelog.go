package actions

import (
	"autogit/actions/changelog"
	"autogit/actions/validation"
	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"autogit/semanticgit/semver"
	"autogit/utils"
)

var CMDChangelog struct {
	VersionParams

	Tag      *string
	Validate *bool
}

func Changelog() string {
	g := (&semanticgit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())
	rendered_changelog := changelog.ChangelogData{Tag: *CMDChangelog.Tag}.New(g, semver.OptionsSemVer{
		DisableVFlag:  *CMDChangelog.DisableVFlag,
		EnableNewline: false,
		Build:         *CMDChangelog.BuildMeta,
		Alpha:         *CMDChangelog.Alpha,
		Beta:          *CMDChangelog.Beta,
		Rc:            *CMDChangelog.Prerelease,
		Publish:       *CMDChangelog.Publish,
	}).Render()

	if *CMDChangelog.Validate {
		log_records := g.GetChangelogByTag(*CMDChangelog.Tag, false)
		for _, record := range log_records {
			err := validation.Validate(&record)
			utils.CheckFatal(err)
		}
	}

	return rendered_changelog
}

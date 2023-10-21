package actions

import (
	"autogit/actions/changelog"
	"autogit/actions/validation"
	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"autogit/settings"
	"autogit/settings/logus"
	"autogit/settings/types"

	"github.com/spf13/cobra"
)

type ChangelogParams struct {
	VersionParams
	Tag      string
	Validate bool
}

func (v *ChangelogParams) Bind(cmd *cobra.Command) {
	v.VersionParams.Bind(cmd)
	cmd.PersistentFlags().StringVar(&v.Tag, "tag", "", "Select from which tag")
	cmd.PersistentFlags().BoolVar(&v.Validate, "validate", false, "Validate to rules")
}

func Changelog(params ChangelogParams, gitw *git.Repository) string {
	conf := settings.GetConfig()
	params.EnableNewline = false

	g := (&semanticgit.SemanticGit{}).NewRepo(gitw)
	rendered_changelog := changelog.ChangelogData{Tag: types.TagName(params.Tag)}.New(g, params.OptionsSemVer).Render()

	if params.Validate {
		log_commits := g.GetChangelogByTag(types.TagName(params.Tag), false)
		for _, commit := range log_commits {
			err := validation.Validate(commit, conf)
			logus.CheckError(err, "failed to validate", logus.Commit(commit.ParsedCommit))
		}
	}

	return rendered_changelog
}

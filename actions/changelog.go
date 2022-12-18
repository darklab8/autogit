package actions

import (
	"autogit/actions/changelog"
	"autogit/actions/validation"
	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"autogit/utils"

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
	params.EnableNewline = false

	g := (&semanticgit.SemanticGit{}).NewRepo(gitw)
	rendered_changelog := changelog.ChangelogData{Tag: params.Tag}.New(g, params.OptionsSemVer).Render()

	if params.Validate {
		log_records := g.GetChangelogByTag(params.Tag, false)
		for _, record := range log_records {
			err := validation.Validate(&record)
			utils.CheckFatal(err)
		}
	}

	return rendered_changelog
}

package actions

import (
	"fmt"

	"github.com/darklab8/autogit/v2/interface_cli/actions/changelog"
	"github.com/darklab8/autogit/v2/interface_cli/actions/validation"
	"github.com/darklab8/autogit/v2/semanticgit"
	"github.com/darklab8/autogit/v2/semanticgit/git"
	"github.com/darklab8/autogit/v2/settings"
	"github.com/darklab8/autogit/v2/settings/logus"
	"github.com/darklab8/autogit/v2/settings/types"

	"github.com/spf13/cobra"
)

type ChangelogParams struct {
	VersionParams
	Tag      string
	Validate bool
	Format   string
}

func (v *ChangelogParams) Bind(cmd *cobra.Command) {
	v.VersionParams.Bind(cmd)
	cmd.PersistentFlags().StringVar(&v.Tag, "tag", "", "Select from which tag")
	cmd.PersistentFlags().BoolVar(&v.Validate, "validate", false, "Validate to rules")
	cmd.PersistentFlags().StringVar(&v.Format, "format", string(changelog.FormatMarkdown), fmt.Sprintf("expected formats=%v", changelog.Formats))
}

func Changelog(params ChangelogParams, gitw *git.Repository) string {
	conf := settings.GetConfig()
	params.EnableNewline = false

	g := semanticgit.NewSemanticRepo(gitw)
	logus.Log.Debug("Getting changelog", logus.TagName(types.TagName(params.Tag)))

	var changelogus changelog.IChangelog
	switch changelog.ChangelogFormat(params.Format) {
	case changelog.FormatMarkdown:
		changelogus = changelog.NewChangelogMarkdown(g, params.OptionsSemVer, conf, types.TagName(params.Tag))
	case changelog.FormatBBCode:
		changelogus = changelog.NewChangelogBBCode(g, params.OptionsSemVer, conf, types.TagName(params.Tag))
	default:
		changelogus = changelog.NewChangelogMarkdown(g, params.OptionsSemVer, conf, types.TagName(params.Tag))
	}
	rendered_changelog := changelogus.Render()

	if params.Validate {
		var anyError error = nil

		log_commits := g.GetChangelogByTag(types.TagName(params.Tag), false)
		for _, commit := range log_commits {
			err := validation.Validate(commit, conf)
			if logus.Log.CheckError(err, "failed to validate", logus.Commit(commit.ParsedCommit)) {
				anyError = err
			}
		}

		logus.Log.CheckFatal(anyError, "encountered at least one error during changelog validation")
	}

	return rendered_changelog
}

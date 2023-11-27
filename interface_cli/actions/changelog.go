package actions

import (
	"autogit/interface_cli/actions/changelog"
	"autogit/interface_cli/actions/validation"
	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"autogit/settings"
	"autogit/settings/logus"
	"autogit/settings/types"
	"fmt"

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
	logus.Debug("Getting changelog", logus.TagName(types.TagName(params.Tag)))

	var changelogus changelog.IChangelog
	switch changelog.ChangelogFormat(params.Format) {
	case changelog.FormatMarkdown:
		changelogus = changelog.NewChangelogMarkdown(g, params.OptionsSemVer, conf, types.TagName(params.Tag))
	case changelog.FormatBBCode:
		changelogus = changelog.NewChangelogBBCode(g, params.OptionsSemVer, conf, types.TagName(params.Tag))
	default:
		logus.Error(fmt.Sprintf("Not allowed format=%s, expected=%v", params.Format, changelog.Formats))
		return ""
	}
	rendered_changelog := changelogus.Render()

	if params.Validate {
		log_commits := g.GetChangelogByTag(types.TagName(params.Tag), false)
		for _, commit := range log_commits {
			err := validation.Validate(commit, conf)
			logus.CheckError(err, "failed to validate", logus.Commit(commit.ParsedCommit))
		}
	}

	return rendered_changelog
}

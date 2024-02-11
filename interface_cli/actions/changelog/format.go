package changelog

import "github.com/darklab8/autogit/interface_cli/actions/changelog/changelog_types"

type IChangelog interface {
	Render() string
	GetSemverGroups() map[changelog_types.ChangelogSectionType]*changelogSemverGroup
}

type ChangelogFormat string

const (
	FormatMarkdown ChangelogFormat = "markdown"
	FormatBBCode   ChangelogFormat = "bbcode"
)

var Formats = []ChangelogFormat{FormatMarkdown, FormatBBCode}

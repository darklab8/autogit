package changelog

import (
	"fmt"
	"strings"

	"github.com/darklab8/autogit/v2/interface_cli/actions/changelog/changelog_types"
	"github.com/darklab8/autogit/v2/interface_cli/actions/changelog/templates"
	"github.com/darklab8/autogit/v2/semanticgit"
	"github.com/darklab8/autogit/v2/semanticgit/conventionalcommits"
	"github.com/darklab8/autogit/v2/semanticgit/conventionalcommits/conventionalcommitstype"
	"github.com/darklab8/autogit/v2/semanticgit/semver/semvertype"
	"github.com/darklab8/autogit/v2/settings"
	"github.com/darklab8/autogit/v2/settings/types"
)

type ChangelogMarkdown struct {
	Changelog
}

func NewChangelogMarkdown(
	g *semanticgit.SemanticGit,
	semver_options semvertype.OptionsSemVer,
	config settings.ConfigScheme,
	FromTag types.TagName,
) IChangelog {
	return ChangelogMarkdown{Changelog: NewChangelog(g, semver_options, config, FromTag)}
}

func (changelog ChangelogMarkdown) GetSemverGroups() map[changelog_types.ChangelogSectionType]*changelogSemverGroup {
	return changelog.SemverGroups
}

func (changelog ChangelogMarkdown) RenderHeader() string {
	return fmt.Sprintf("## **%s** <sub><sub>%s ([%s...%s](%s))</sub></sub>",
		changelog.Header.ChangelogVersion,
		changelog.Header.Timestamp,
		changelog.Header.From,
		changelog.Header.To,
		changelog.Header.CommitRangeURL,
	)
}

func (changelog ChangelogMarkdown) commitRender(record conventionalcommits.ConventionalCommit) changelog_types.ChangelogCommitHeader {
	templs := templates.NewTemplates()

	issue_renderer := strings.Builder{}
	for _, issue_n := range record.Issue {
		issue_renderer.WriteString(fmt.Sprintf(", [#%s](%s)", issue_n, templs.RenderIssueUrl(issue_n)))
	}
	issue_rendered := issue_renderer.String()

	rendered_subject := record.Subject
	IssueMatch := settings.RegexIssue.FindAllStringSubmatch(string(record.Subject), -1)
	for _, match := range IssueMatch {
		rendered_subject = conventionalcommitstype.Subject(strings.Replace(
			string(rendered_subject),
			match[0],
			fmt.Sprintf(
				"[#%s](%s)",
				match[1],
				templs.RenderIssueUrl(conventionalcommitstype.Issue(match[1])),
			),
			-1,
		),
		)
	}

	formatted_url := templs.RenderCommitUrl(record)
	result := fmt.Sprintf("* %s ([%s](%s)%s)\n", rendered_subject, record.Hash, formatted_url, issue_rendered)
	return changelog_types.ChangelogCommitHeader(result)
}

func (changelog ChangelogMarkdown) Render() string {
	var sb strings.Builder = strings.Builder{}
	sbprintln := func(format string, a ...any) {
		sb.WriteString(fmt.Sprintf(format+"\n", a...))
	}
	sbprint := func(format string, a ...any) {
		sb.WriteString(fmt.Sprintf(format, a...))
	}

	print_commits := func(prefix string, commits []changelogCommit) {
		for _, commit := range commits {
			sbprint(prefix+"%s", changelog.commitRender(commit.Commit))

			for _, breaking_footer := range commit.BreakingFooters {
				breaking_footer_lines := strings.Split(string(breaking_footer), "\n")
				sbprintln(prefix+"  * BC!: %s", breaking_footer_lines[0])
				for _, other_breaking_change_line := range breaking_footer_lines[1:] {
					sbprintln(prefix+"        %s", other_breaking_change_line)
				}
			}
		}
	}

	sbprintln("%s", changelog.RenderHeader())
	sbprintln("")
	for _, semver_group := range changelog.OrderedSemverGroups {
		sbprintln("## %s", semver_group.Name)
		for commit_type, type_group := range semver_group.CommitTypeGroups {
			sbprintln("### %s", commit_type)

			print_commits("", type_group.NoScopeCommits)
			for scope, commits := range type_group.ScopedCommits {
				sbprintln("* %s", scope)
				print_commits("  ", commits)
			}
		}
		sbprintln("")
	}

	return sb.String()
}

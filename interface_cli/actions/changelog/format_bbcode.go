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

type ChangelogBBCode struct {
	Changelog
}

func NewChangelogBBCode(
	g *semanticgit.SemanticGit,
	semver_options semvertype.OptionsSemVer,
	config settings.ConfigScheme,
	FromTag types.TagName,
) IChangelog {
	return ChangelogBBCode{Changelog: NewChangelog(g, semver_options, config, FromTag)}
}

func (changelog ChangelogBBCode) GetSemverGroups() map[changelog_types.ChangelogSectionType]*changelogSemverGroup {
	return changelog.SemverGroups
}

func (changelog ChangelogBBCode) RenderHeader() string {
	return fmt.Sprintf("[size=large][color=#00FF00][b]%s[/b][/color] [color=#FFFF00]%s[/color] [url=%s](%s...%s)[/url][/size]",
		changelog.Header.ChangelogVersion,
		changelog.Header.Timestamp,
		changelog.Header.CommitRangeURL,
		changelog.Header.From,
		changelog.Header.To,
	)
}

func (changelog ChangelogBBCode) commitRender(record conventionalcommits.ConventionalCommit) changelog_types.ChangelogCommitHeader {
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
				"[url=%s](%s)[/url]",
				match[1],
				templs.RenderIssueUrl(conventionalcommitstype.Issue(match[1])),
			),
			-1,
		),
		)
	}

	formatted_url := templs.RenderCommitUrl(record)
	result := fmt.Sprintf("[*] %s [url=%s](%s)[/url] %s\n", rendered_subject, formatted_url, record.Hash, issue_rendered)
	return changelog_types.ChangelogCommitHeader(result)
}

func (changelog ChangelogBBCode) Render() string {
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

			if len(commit.BreakingFooters) > 0 {
				sbprintln("  [list]")
				for _, breaking_footer := range commit.BreakingFooters {
					breaking_footer_lines := strings.Split(string(breaking_footer), "\n")
					sbprintln(prefix+"  [*] BC!: %s", breaking_footer_lines[0])
					for _, other_breaking_change_line := range breaking_footer_lines[1:] {
						sbprintln(prefix+"        %s", other_breaking_change_line)
					}
				}
				sbprintln("  [/list]")
			}

		}
	}

	sbprintln("[quote=changelog]")

	sbprintln("%s", changelog.RenderHeader())
	for _, semver_group := range changelog.OrderedSemverGroups {
		sbprintln("[size=large]%s[/size]", semver_group.Name)
		for commit_type, type_group := range semver_group.CommitTypeGroups {
			sbprintln("[size=medium]%s[/size]", commit_type)

			if len(type_group.NoScopeCommits) > 0 || len(type_group.ScopedCommits) > 0 {

				sbprintln("[list]")
				print_commits("", type_group.NoScopeCommits)

				for scope, commits := range type_group.ScopedCommits {
					sbprint("[*] %s ", scope)
					if len(commits) > 0 {
						sbprintln("[list]")
						print_commits("  ", commits)
						sbprintln("  [/list]")
					}
				}

				sbprintln("[/list]")
			}
		}
		sbprintln("")
	}
	sbprintln("[/quote]")

	return sb.String()
}

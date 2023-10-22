package changelog

import (
	"autogit/interface_cli/actions/changelog/templates"
	"autogit/semanticgit"
	"autogit/semanticgit/conventionalcommits"
	"autogit/semanticgit/conventionalcommits/conventionalcommitstype"
	"autogit/semanticgit/semver/semvertype"
	"autogit/settings"
	"autogit/settings/logus"
	"autogit/settings/types"
	"autogit/settings/utils"
	"fmt"
	"strings"
	"text/template"

	_ "embed"
)

func commitRender(record conventionalcommits.ConventionalCommit) string {
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
	return result
}

type commitTypeGroup struct {
	NoScopeCommits []string
	ScopedCommits  map[conventionalcommitstype.Scope][]string
}

type changelogSemverGroup struct {
	CommitTypeGroups map[conventionalcommitstype.Type]*commitTypeGroup
	Name             SemverOrder
}

type SemverOrder string

const (
	Major SemverOrder = "Major"
	Minor SemverOrder = "Minor"
	Patch SemverOrder = "Patch"
)

type changelogVars struct {
	// Internal for data grouping
	SemverGroups map[SemverOrder]*changelogSemverGroup

	// For template
	Header              string
	OrderedSemverGroups []*changelogSemverGroup
}

func (changelog *changelogVars) find_semver_group(
	record conventionalcommits.ConventionalCommit,
	types []conventionalcommitstype.Type,
	semver_order SemverOrder,
) (*changelogSemverGroup, error) {
	for _, possible_type := range types {
		if record.Exclamation {
			semver_group, semver_group_exists := changelog.SemverGroups[Major]
			if !semver_group_exists {
				semver_group = &changelogSemverGroup{Name: Major}
				changelog.SemverGroups[Major] = semver_group
			}
			return semver_group, nil
		}

		if record.Type == possible_type {
			semver_group, semver_group_exists := changelog.SemverGroups[semver_order]
			if !semver_group_exists {
				semver_group = &changelogSemverGroup{Name: semver_order}
				changelog.SemverGroups[semver_order] = semver_group
			}
			return semver_group, nil
		}
	}
	return nil, NotFound{}
}

type NotFound struct {
}

func (n NotFound) Error() string {
	return "not found SemverGroup"
}

func (changelog *changelogVars) addCommit(
	record conventionalcommits.ConventionalCommit,
	commit_formatted string,
	config settings.ConfigScheme,
) {

	semver_group, err := changelog.find_semver_group(record, config.Validation.Rules.Header.Type.Allowlists.SemVerMinorIncreasers, Minor)
	if err != nil {
		semver_group, err = changelog.find_semver_group(record, config.Validation.Rules.Header.Type.Allowlists.SemverPatchIncreasers, Patch)
	}
	if err != nil {
		return
	}

	if semver_group.CommitTypeGroups == nil {
		semver_group.CommitTypeGroups = make(map[conventionalcommitstype.Type]*commitTypeGroup)
	}

	commit_type_group, ok := semver_group.CommitTypeGroups[record.Type]
	if !ok {
		commit_type_group = &commitTypeGroup{}
		semver_group.CommitTypeGroups[record.Type] = commit_type_group
	}

	if record.Scope == "" {
		commit_type_group.NoScopeCommits = append(commit_type_group.NoScopeCommits, commit_formatted)
	} else {
		if commit_type_group.ScopedCommits == nil {
			commit_type_group.ScopedCommits = make(map[conventionalcommitstype.Scope][]string)
		}

		commit_list, ok := commit_type_group.ScopedCommits[record.Scope]
		if !ok {
			commit_list = []string{}
		}

		commit_type_group.ScopedCommits[record.Scope] = append(commit_list, commit_formatted)
	}
}

func NewChangelog(g *semanticgit.SemanticGit, semver_options semvertype.OptionsSemVer, config settings.ConfigScheme, FromTag types.TagName) changelogVars {
	templs := templates.NewTemplates()

	changelog := changelogVars{}
	changelog.SemverGroups = make(map[SemverOrder]*changelogSemverGroup)

	logs := g.GetChangelogByTag(FromTag, true)
	logus.Debug(fmt.Sprintf("NewChangelog, log.count=%d", len(logs)))
	if FromTag == "" {
		FromTag = g.GetNextVersion(semver_options).ToString()
	}

	ChangelogVersionTag := FromTag
	if FromTag == "" {
		ChangelogVersionTag = g.GetNextVersion(semver_options).ToString()
	}

	changelog.Header = templs.NewCommitRangeUrlRender(logs, ChangelogVersionTag)

	for _, record := range logs {
		commit_formatted := commitRender(record)
		changelog.addCommit(record, commit_formatted, config)
	}

	// for easier templating as ordered
	if semver_group, found := changelog.SemverGroups[Major]; found {
		changelog.OrderedSemverGroups = append(changelog.OrderedSemverGroups, semver_group)
	}
	if semver_group, found := changelog.SemverGroups[Minor]; found {
		changelog.OrderedSemverGroups = append(changelog.OrderedSemverGroups, semver_group)
	}
	if semver_group, found := changelog.SemverGroups[Patch]; found {
		changelog.OrderedSemverGroups = append(changelog.OrderedSemverGroups, semver_group)
	}

	return changelog
}

func (changelog changelogVars) Render() string {
	return utils.TmpRender(changelogTemplate, changelog)
}

//go:embed templates/changelog.md
var changelogMarkup types.TemplateExpression
var changelogTemplate *template.Template

func init() {
	changelogTemplate = utils.TmpInit(changelogMarkup)
}

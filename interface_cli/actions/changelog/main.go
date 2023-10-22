package changelog

import (
	"autogit/interface_cli/actions/changelog/templates"
	"autogit/semanticgit"
	"autogit/semanticgit/conventionalcommits"
	"autogit/semanticgit/conventionalcommits/conventionalcommitstype"
	"autogit/semanticgit/semver"
	"autogit/settings"
	"autogit/settings/types"
	"autogit/settings/utils"
	"fmt"
	"strings"
	"text/template"

	_ "embed"
)

func commitRender(record conventionalcommits.ConventionalCommit) string {
	templs := templates.NewTemplates()

	var issue_rendered strings.Builder
	for _, issue_n := range record.Issue {

		issue_rendered.WriteString(fmt.Sprintf(", [#%s](%s)", issue_n, templs.RenderIssueUrl(issue_n)))
	}

	rendered_subject := record.Subject
	IssueMatch := settings.RegexIssue.FindAllStringSubmatch(record.Subject, -1)
	for _, match := range IssueMatch {
		rendered_subject = strings.Replace(rendered_subject, match[0], fmt.Sprintf("[#%s](%s)", match[1], templs.RenderIssueUrl(match[1])), -1)
	}

	formatted_url := templs.RenderCommitUrl(record)
	return fmt.Sprintf("* %s ([%s](%s)%s)\n", rendered_subject, record.Hash, formatted_url, issue_rendered.String())
}

type changelogVars struct {
	Header           string
	Features         []string
	Fixes            []string
	FeaturesScoped   map[conventionalcommitstype.Scope][]string
	FixesScoped      map[conventionalcommitstype.Scope][]string
	AreThereFeatures bool
	AreThereFixes    bool
}

func (changelog *changelogVars) addCommit(record conventionalcommits.ConventionalCommit, commit_formatted string) {
	if record.Type == "feat" {
		if record.Scope == "" {
			changelog.Features = append(changelog.Features, commit_formatted)
		} else {
			_, ok := changelog.FeaturesScoped[record.Scope]
			if !ok {
				changelog.FeaturesScoped[record.Scope] = []string{}
			}

			changelog.FeaturesScoped[record.Scope] = append(changelog.FeaturesScoped[record.Scope], commit_formatted)
		}
	} else if record.Type == "fix" {
		if record.Scope == "" {
			changelog.Fixes = append(changelog.Fixes, commit_formatted)
		} else {
			_, ok := changelog.FixesScoped[record.Scope]
			if !ok {
				changelog.FixesScoped[record.Scope] = []string{}
			}

			changelog.FixesScoped[record.Scope] = append(changelog.FixesScoped[record.Scope], commit_formatted)
		}
	}
}

func NewChangelog(g *semanticgit.SemanticGit, semver_options semver.OptionsSemVer, config settings.ChangelogScheme, FromTag types.TagName) changelogVars {
	changelog := changelogVars{}
	changelog.FeaturesScoped = make(map[conventionalcommitstype.Scope][]string)
	changelog.FixesScoped = make(map[conventionalcommitstype.Scope][]string)

	templs := templates.NewTemplates()

	logs := g.GetChangelogByTag(FromTag, true)
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
		changelog.addCommit(record, commit_formatted)
	}

	changelog.AreThereFeatures = len(changelog.Features) > 0 || len(changelog.FeaturesScoped) > 0
	changelog.AreThereFixes = len(changelog.Fixes) > 0 || len(changelog.FixesScoped) > 0

	return changelog
}

func (changelog changelogVars) Render() string {
	return utils.TmpRender(changelogTemplate, changelog)
}

//go:embed templates/changelog.md
var changelogMarkup string
var changelogTemplate *template.Template

func init() {
	changelogTemplate = utils.TmpInit(changelogMarkup)
}

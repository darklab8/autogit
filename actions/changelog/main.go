package changelog

import (
	"autogit/semanticgit"
	"autogit/semanticgit/conventionalcommits"
	"autogit/semanticgit/semver"
	"autogit/settings"
	"autogit/utils"
	"fmt"
	"strings"
	"text/template"
	"time"

	_ "embed"
)

type Header struct {
	From    string
	To      string
	Version string
}

func (r *Header) Render() string {
	templs := settings.GetTemplates()
	currentTime := time.Now()
	return fmt.Sprintf("## **%s** <sub><sub>%s ([%s...%s](%s))</sub></sub>", r.Version, currentTime.Format("2006-01-02"), r.From, r.To, utils.TmpRender(templs.CommitRangeUrl, r))
}

func (r *Header) New(logs []conventionalcommits.ConventionalCommit, version string) *Header {
	r.From = logs[len(logs)-1].Hash
	r.To = logs[0].Hash
	r.Version = version
	return r
}

type commitRecord struct {
	Commit string
}

func (c commitRecord) Render(record conventionalcommits.ConventionalCommit) string {
	templs := settings.GetTemplates()
	type IssueData struct {
		Issue string
	}

	var issue_rendered strings.Builder
	for _, issue_n := range record.Issue {
		issue_rendered.WriteString(fmt.Sprintf(", [#%s](%s)", issue_n, utils.TmpRender(templs.IssueUrl, IssueData{Issue: issue_n})))
	}

	rendered_subject := record.Subject
	IssueMatch := settings.RegexIssue.FindAllStringSubmatch(record.Subject, -1)
	for _, match := range IssueMatch {
		rendered_subject = strings.Replace(rendered_subject, match[0], fmt.Sprintf("[#%s](%s)", match[1], utils.TmpRender(templs.IssueUrl, IssueData{Issue: match[1]})), -1)
	}

	formatted_url := utils.TmpRender(templs.CommitUrl, commitRecord{Commit: record.Hash})
	return fmt.Sprintf("* %s ([%s](%s)%s)\n", rendered_subject, record.Hash, formatted_url, issue_rendered.String())
}

type ChangelogData struct {
	Tag              string // get changelog from this tag to previous
	Header           string
	Features         []string
	Fixes            []string
	FeaturesScoped   map[string][]string
	FixesScoped      map[string][]string
	AreThereFeatures bool
	AreThereFixes    bool
}

func (changelog *ChangelogData) AddCommit(record conventionalcommits.ConventionalCommit, commit_formatted string) {
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

func (changelog ChangelogData) New(g *semanticgit.SemanticGit, semver_options semver.OptionsSemVer) ChangelogData {
	changelog.FeaturesScoped = make(map[string][]string)
	changelog.FixesScoped = make(map[string][]string)

	logs := g.GetChangelogByTag(changelog.Tag, true)

	if changelog.Tag == "" {
		changelog.Tag = g.GetNextVersion(semver_options).ToString()
	}
	changelog.Header = (&Header{}).New(logs, changelog.Tag).Render()

	for _, record := range logs {
		commit_formatted := commitRecord{Commit: record.Hash}.Render(record)
		changelog.AddCommit(record, commit_formatted)
	}

	changelog.AreThereFeatures = len(changelog.Features) > 0 || len(changelog.FeaturesScoped) > 0
	changelog.AreThereFixes = len(changelog.Fixes) > 0 || len(changelog.FixesScoped) > 0

	return changelog
}

func (changelog ChangelogData) Render() string {
	return utils.TmpRender(changelogTemplate, changelog)
}

//go:embed templates/changelog.md
var changelogMarkup string
var changelogTemplate *template.Template

func init() {
	changelogTemplate = utils.TmpInit(changelogMarkup)
}

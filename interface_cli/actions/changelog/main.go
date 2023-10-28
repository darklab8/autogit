package changelog

import (
	"autogit/interface_cli/actions/changelog/changelog_types"
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
	Name             changelog_types.ChangelogSectionName
}

const (
	SemVerMajor  changelog_types.ChangelogSection = "semver_major"
	SemVerMinor  changelog_types.ChangelogSection = "semver_minor"
	SemVerPatch  changelog_types.ChangelogSection = "semver_patch"
	MergeCommits changelog_types.ChangelogSection = "merge_commits"
)

func GetSectionName(section changelog_types.ChangelogSection) changelog_types.ChangelogSectionName {
	config := settings.GetConfig()
	is_pr := config.Changelog.MergeCommits.MustHaveLinkedPR

	var merge_heading_prefix changelog_types.ChangelogSectionName
	switch section {
	case MergeCommits:
		if is_pr {
			merge_heading_prefix = config.Changelog.Headings.MergeCommits.WithLinkedPR
		} else {
			merge_heading_prefix = config.Changelog.Headings.MergeCommits.Default
		}

		if config.Changelog.MergeCommits.RedirectMergingCommits {
			temp := config.Changelog.Headings.MergeCommits.PrefixForUndirected + " " + string(merge_heading_prefix)
			merge_heading_prefix = changelog_types.ChangelogSectionName(temp)
		}

		return merge_heading_prefix
	case SemVerMajor:
		return config.Changelog.Headings.SemverMajor
	case SemVerMinor:
		return config.Changelog.Headings.SemverMinor
	case SemVerPatch:
		return config.Changelog.Headings.SemverPatch
	default:
		panic("GetSectionName encountered not supported section")
	}
}

type changelogVars struct {
	// Internal for data grouping
	SemverGroups map[changelog_types.ChangelogSection]*changelogSemverGroup

	// For template
	Header              string
	OrderedSemverGroups []*changelogSemverGroup
}

func (changelog *changelogVars) find_semver_group(
	record conventionalcommits.ConventionalCommit,
	types []conventionalcommitstype.Type,
	semver_order changelog_types.ChangelogSection,
) (*changelogSemverGroup, error) {
	for _, possible_type := range types {
		if record.Exclamation {
			semver_group, semver_group_exists := changelog.SemverGroups[SemVerMajor]
			if !semver_group_exists {
				semver_group = &changelogSemverGroup{Name: GetSectionName(SemVerMajor)}
				changelog.SemverGroups[SemVerMajor] = semver_group
			}
			return semver_group, nil
		}

		if record.Type == possible_type {
			semver_group, semver_group_exists := changelog.SemverGroups[semver_order]
			if !semver_group_exists {
				semver_group = &changelogSemverGroup{Name: GetSectionName(semver_order)}
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
	if config.Changelog.MergeCommits.MustHaveLinkedPR {
		for _, merge_type := range config.Changelog.MergeCommits.MergeTypes {
			if record.Type == merge_type {
				match_pr := settings.RegexPullRequest.FindStringSubmatch(string(record.Subject))
				if len(match_pr) == 0 {
					logus.Debug("merging commit is not containing linked PR", logus.Commit(record.ParsedCommit))
					return
				} else {
					logus.Debug(fmt.Sprintf("merging commit contains linked PR=%v", match_pr), logus.Commit(record.ParsedCommit))
				}
			}
		}
	}

	redirect_merging_commit := func(commit *conventionalcommits.ConventionalCommit, redirecting_type conventionalcommitstype.Type) {
		commit.Type = redirecting_type
		// Could be replaced with regex. Or it is fine as it is.
		if strings.Contains(string(commit.Body), fmt.Sprintf("%s!", redirecting_type)) {
			commit.Exclamation = true
		}
		if strings.Contains(string(commit.Body), ")!") {
			commit.Exclamation = true
		}
	}
	redirect := func(commit *conventionalcommits.ConventionalCommit, iterable_types []conventionalcommitstype.Type) error {
		for _, iterated_type := range iterable_types {
			logus.Debug(fmt.Sprintf("RedirectMergingCommits... for type=%s", iterated_type), logus.Commit(commit.ParsedCommit))
			if strings.Contains(string(commit.Subject), string(iterated_type)) {
				redirect_merging_commit(commit, iterated_type)
				return nil
			}
			for _, footer_stuff := range commit.Footers {
				if strings.Contains(string(footer_stuff.Token), string(iterated_type)) {
					redirect_merging_commit(commit, iterated_type)
					return nil
				}
			}
		}
		return NotFound{}
	}
	if config.Changelog.MergeCommits.RedirectMergingCommits {
		for _, merge_type := range config.Changelog.MergeCommits.MergeTypes {
			if record.Type == merge_type {
				if redirect(&record, config.Validation.Rules.Header.Type.Allowlists.SemVerMinorIncreasers) == nil {
					break
				}
				if redirect(&record, config.Validation.Rules.Header.Type.Allowlists.SemverPatchIncreasers) == nil {
					break
				}
			}
		}
	}

	semver_group, err := changelog.find_semver_group(record, config.Validation.Rules.Header.Type.Allowlists.SemVerMinorIncreasers, SemVerMinor)
	if err != nil {
		semver_group, err = changelog.find_semver_group(record, config.Validation.Rules.Header.Type.Allowlists.SemverPatchIncreasers, SemVerPatch)
	}
	if err != nil {
		semver_group, err = changelog.find_semver_group(record, config.Changelog.MergeCommits.MergeTypes, MergeCommits)
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
	changelog.SemverGroups = make(map[changelog_types.ChangelogSection]*changelogSemverGroup)

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
	if semver_group, found := changelog.SemverGroups[MergeCommits]; found {
		changelog.OrderedSemverGroups = append(changelog.OrderedSemverGroups, semver_group)
	}
	if semver_group, found := changelog.SemverGroups[SemVerMajor]; found {
		changelog.OrderedSemverGroups = append(changelog.OrderedSemverGroups, semver_group)
	}
	if semver_group, found := changelog.SemverGroups[SemVerMinor]; found {
		changelog.OrderedSemverGroups = append(changelog.OrderedSemverGroups, semver_group)
	}
	if semver_group, found := changelog.SemverGroups[SemVerPatch]; found {
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

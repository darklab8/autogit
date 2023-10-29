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
	"fmt"
	"strings"
)

func commitRender(record conventionalcommits.ConventionalCommit) changelog_types.ChangelogCommitHeader {
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

type changelogCommit struct {
	Header          changelog_types.ChangelogCommitHeader
	BreakingFooters []conventionalcommitstype.FooterContent
}

func newChangelogCommit(record conventionalcommits.ConventionalCommit) changelogCommit {
	header := commitRender(record)
	changelog_commit := changelogCommit{Header: header}

	for _, footer := range record.Footers {
		if footer.Token == semanticgit.FooterTokenBreakingChange {
			changelog_commit.BreakingFooters = append(changelog_commit.BreakingFooters, footer.Content)
		}
	}

	return changelog_commit
}

type commitTypeGroup struct {
	NoScopeCommits []changelogCommit
	ScopedCommits  map[conventionalcommitstype.Scope][]changelogCommit
}

type changelogSemverGroup struct {
	CommitTypeGroups map[conventionalcommitstype.Type]*commitTypeGroup
	Name             changelog_types.ChangelogSectionName
}

func getSectionName(section changelog_types.ChangelogSectionType) changelog_types.ChangelogSectionName {
	config := settings.GetConfig()
	is_pr := config.Changelog.MergeCommits.MustHaveLinkedPR

	var merge_heading_prefix changelog_types.ChangelogSectionName
	switch section {
	case changelog_types.MergeCommits:
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
	case changelog_types.SemVerMajor:
		return config.Changelog.Headings.SemverMajor
	case changelog_types.SemVerMinor:
		return config.Changelog.Headings.SemverMinor
	case changelog_types.SemVerPatch:
		return config.Changelog.Headings.SemverPatch
	default:
		panic("GetSectionName encountered not supported section")
	}
}

type changelogVars struct {
	// Internal for data grouping
	SemverGroups map[changelog_types.ChangelogSectionType]*changelogSemverGroup

	// For template
	Header              string
	OrderedSemverGroups []*changelogSemverGroup
}

func (changelog *changelogVars) find_semver_group(
	record conventionalcommits.ConventionalCommit,
	conventiona_types []conventionalcommitstype.Type,
	semver_order changelog_types.ChangelogSectionType,
) (*changelogSemverGroup, error) {
	for _, possible_type := range conventiona_types {
		if semanticgit.IsBreakingChangeCommit(record) {
			semver_group, semver_group_exists := changelog.SemverGroups[changelog_types.SemVerMajor]
			if !semver_group_exists {
				semver_group = &changelogSemverGroup{Name: getSectionName(changelog_types.SemVerMajor)}
				changelog.SemverGroups[changelog_types.SemVerMajor] = semver_group
			}
			return semver_group, nil
		}

		if record.Type == possible_type {
			semver_group, semver_group_exists := changelog.SemverGroups[semver_order]
			if !semver_group_exists {
				semver_group = &changelogSemverGroup{Name: getSectionName(semver_order)}
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
		// P.S. Not very reliable mechanism to identify breaking change commits in merging commits :/
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

	semver_group, err := changelog.find_semver_group(record, config.Validation.Rules.Header.Type.Allowlists.SemVerMinorIncreasers, changelog_types.SemVerMinor)
	if err != nil {
		semver_group, err = changelog.find_semver_group(record, config.Validation.Rules.Header.Type.Allowlists.SemverPatchIncreasers, changelog_types.SemVerPatch)
	}
	if err != nil {
		semver_group, err = changelog.find_semver_group(record, config.Changelog.MergeCommits.MergeTypes, changelog_types.MergeCommits)
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

	changelog_commit := newChangelogCommit(record)

	if record.Scope == "" {
		commit_type_group.NoScopeCommits = append(commit_type_group.NoScopeCommits, changelog_commit)
	} else {
		if commit_type_group.ScopedCommits == nil {
			commit_type_group.ScopedCommits = make(map[conventionalcommitstype.Scope][]changelogCommit)
		}

		commit_list, ok := commit_type_group.ScopedCommits[record.Scope]
		if !ok {
			commit_list = []changelogCommit{}
		}

		commit_type_group.ScopedCommits[record.Scope] = append(commit_list, changelog_commit)
	}
}

func (changelog changelogVars) orderSemverGroups() []*changelogSemverGroup {
	result := []*changelogSemverGroup{}

	if semver_group, found := changelog.SemverGroups[changelog_types.MergeCommits]; found {
		result = append(result, semver_group)
	}
	if semver_group, found := changelog.SemverGroups[changelog_types.SemVerMajor]; found {
		result = append(result, semver_group)
	}
	if semver_group, found := changelog.SemverGroups[changelog_types.SemVerMinor]; found {
		result = append(result, semver_group)
	}
	if semver_group, found := changelog.SemverGroups[changelog_types.SemVerPatch]; found {
		result = append(result, semver_group)
	}
	return result
}

func NewChangelog(g *semanticgit.SemanticGit, semver_options semvertype.OptionsSemVer, config settings.ConfigScheme, FromTag types.TagName) changelogVars {
	templs := templates.NewTemplates()

	changelog := changelogVars{}
	changelog.SemverGroups = make(map[changelog_types.ChangelogSectionType]*changelogSemverGroup)

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
		changelog.addCommit(record, config)
	}

	changelog.OrderedSemverGroups = changelog.orderSemverGroups()
	return changelog
}

func (changelog changelogVars) Render() string {
	var sb strings.Builder = strings.Builder{}
	sbprintln := func(format string, a ...any) {
		sb.WriteString(fmt.Sprintf(format+"\n", a...))
	}
	sbprint := func(format string, a ...any) {
		sb.WriteString(fmt.Sprintf(format, a...))
	}

	print_commits := func(prefix string, commits []changelogCommit) {
		for _, commit := range commits {
			sbprint(prefix+"%s", commit.Header)

			for _, breaking_footer := range commit.BreakingFooters {
				breaking_footer_lines := strings.Split(string(breaking_footer), "\n")
				sbprintln(prefix+"  * BC!: %s", breaking_footer_lines[0])
				for _, other_breaking_change_line := range breaking_footer_lines[1:] {
					sbprintln(prefix+"        %s", other_breaking_change_line)
				}
			}
		}
	}

	sbprintln("%s", changelog.Header)
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

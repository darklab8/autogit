package settings

import (
	"github.com/darklab8/autogit/v2/interface_cli/actions/changelog/changelog_types"
	"github.com/darklab8/autogit/v2/semanticgit/conventionalcommits/conventionalcommitstype"
	"github.com/darklab8/autogit/v2/settings/logus"

	"github.com/darklab8/go-utils/utils/utils_types"
)

type ChangelogScheme struct {
	REPOSITORY_OWNER string                         `yaml:"REPOSITORY_OWNER"`
	REPOSITORY_NAME  string                         `yaml:"REPOSITORY_NAME"`
	CommitURL        utils_types.TemplateExpression `yaml:"commitUrl"`
	CommitRangeURL   utils_types.TemplateExpression `yaml:"commitRangeUrl"`
	IssueURL         utils_types.TemplateExpression `yaml:"issueUrl"`

	MergeCommits struct {
		MustHaveLinkedPR       bool                           `yaml:"must_have_linked_pull_request"`
		RedirectMergingCommits bool                           `yaml:"redirect_merging_to_semver_sections_for_changelog"`
		MergeTypes             []conventionalcommitstype.Type `yaml:"commit_types"`
	} `yaml:"merge_commits"`

	Headings struct {
		SemverMajor  changelog_types.ChangelogSectionName `yaml:"semver_major"`
		SemverMinor  changelog_types.ChangelogSectionName `yaml:"semver_minor"`
		SemverPatch  changelog_types.ChangelogSectionName `yaml:"semver_patch"`
		MergeCommits struct {
			Default             changelog_types.ChangelogSectionName `yaml:"default"`
			WithLinkedPR        changelog_types.ChangelogSectionName `yaml:"with_linked_pr"`
			PrefixForUndirected string                               `yaml:"prefix_for_undirected"`
		} `yaml:"merge_commits"`
	}
}

func (conf ConfigScheme) changelogValidate() {
	if conf.Changelog.CommitURL == "" {
		logus.Log.Fatal("autogit.yml->Changelog.CommitUrl is empty")
	}
	if conf.Changelog.CommitRangeURL == "" {
		logus.Log.Fatal("autogit.yml->Changelog.CommitRangeURL is empty")
	}
	if conf.Changelog.IssueURL == "" {
		logus.Log.Fatal("autogit.yml->Changelog.IssueURL is empty")
	}
}

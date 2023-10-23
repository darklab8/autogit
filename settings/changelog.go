package settings

import (
	"autogit/semanticgit/conventionalcommits/conventionalcommitstype"
	"autogit/settings/logus"
	"autogit/settings/types"
)

type ChangelogScheme struct {
	REPOSITORY_OWNER string                   `yaml:"REPOSITORY_OWNER"`
	REPOSITORY_NAME  string                   `yaml:"REPOSITORY_NAME"`
	CommitURL        types.TemplateExpression `yaml:"commitUrl"`
	CommitRangeURL   types.TemplateExpression `yaml:"commitRangeUrl"`
	IssueURL         types.TemplateExpression `yaml:"issueUrl"`

	MergeCommits struct {
		MustHaveLinkedPR       bool                           `yaml:"must_have_linked_pull_request"`
		RedirectMergingCommits bool                           `yaml:"redirect_merging_to_semver_sections_for_changelog"`
		MergeTypes             []conventionalcommitstype.Type `yaml:"commit_types"`
	} `yaml:"merge_commits"`
}

func (conf ConfigScheme) changelogValidate() {
	if conf.Changelog.CommitURL == "" {
		logus.Fatal("autogit.yml->Changelog.CommitUrl is empty")
	}
	if conf.Changelog.CommitRangeURL == "" {
		logus.Fatal("autogit.yml->Changelog.CommitRangeURL is empty")
	}
	if conf.Changelog.IssueURL == "" {
		logus.Fatal("autogit.yml->Changelog.IssueURL is empty")
	}
}

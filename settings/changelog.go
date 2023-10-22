package settings

import (
	"autogit/settings/logus"
	"autogit/settings/types"
)

type ChangelogScheme struct {
	REPOSITORY_OWNER string                   `yaml:"REPOSITORY_OWNER"`
	REPOSITORY_NAME  string                   `yaml:"REPOSITORY_NAME"`
	CommitURL        types.TemplateExpression `yaml:"commitUrl"`
	CommitRangeURL   types.TemplateExpression `yaml:"commitRangeUrl"`
	IssueURL         types.TemplateExpression `yaml:"issueUrl"`
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

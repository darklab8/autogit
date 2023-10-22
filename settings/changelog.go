package settings

import (
	"autogit/settings/logus"
)

type ChangelogScheme struct {
	REPOSITORY_OWNER string `yaml:"REPOSITORY_OWNER"`
	REPOSITORY_NAME  string `yaml:"REPOSITORY_NAME"`
	CommitURL        string `yaml:"commitUrl"`
	CommitRangeURL   string `yaml:"commitRangeUrl"`
	IssueURL         string `yaml:"issueUrl"`
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

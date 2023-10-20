package settings

import (
	"autogit/settings/logus"
)

type ChangelogScheme struct {
	CommitURL      string `yaml:"commitUrl"`
	CommitRangeURL string `yaml:"commitRangeUrl"`
	IssueURL       string `yaml:"issueUrl"`
}

func ChangelogValidate(conf ConfigScheme) {
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

func ChangelogInit(conf ConfigScheme) {
	ChangelogValidate(conf)
}

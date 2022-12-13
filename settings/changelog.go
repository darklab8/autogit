package settings

import (
	_ "embed"
	"log"
)

type ChangelogScheme struct {
	CommitURL      string `yaml:"commitUrl"`
	CommitRangeURL string `yaml:"commitRangeUrl"`
	IssueURL       string `yaml:"issueUrl"`
}

func ChangelogValidate() {
	if Config.Changelog.CommitURL == "" {
		log.Fatal("autogit.yml->Changelog.CommitUrl is empty")
	}
	if Config.Changelog.CommitRangeURL == "" {
		log.Fatal("autogit.yml->Changelog.CommitRangeURL is empty")
	}
	if Config.Changelog.IssueURL == "" {
		log.Fatal("autogit.yml->Changelog.IssueURL is empty")
	}
}

func ChangelogInit() {
	ChangelogValidate()
}

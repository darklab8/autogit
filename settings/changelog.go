package settings

import (
	"log"
)

type ChangelogScheme struct {
	CommitURL      string `yaml:"commitUrl"`
	CommitRangeURL string `yaml:"commitRangeUrl"`
	IssueURL       string `yaml:"issueUrl"`
}

func ChangelogValidate(conf ConfigScheme) {
	if conf.Changelog.CommitURL == "" {
		log.Fatal("autogit.yml->Changelog.CommitUrl is empty")
	}
	if conf.Changelog.CommitRangeURL == "" {
		log.Fatal("autogit.yml->Changelog.CommitRangeURL is empty")
	}
	if conf.Changelog.IssueURL == "" {
		log.Fatal("autogit.yml->Changelog.IssueURL is empty")
	}
}

func ChangelogInit(conf ConfigScheme) {
	ChangelogValidate(conf)
}

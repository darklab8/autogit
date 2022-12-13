package settings

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type ChangelogScheme struct {
	Changelog struct {
		CommitURL      string `yaml:"commitUrl"`
		CommitRangeURL string `yaml:"commitRangeUrl"`
		IssueURL       string `yaml:"issueUrl"`
	} `yaml:"changelog"`
	Regex struct {
		ConventionalCommit struct {
			Header     string `yaml:"header"`
			BodyFooter string `yaml:"bodyFooter"`
		} `yaml:"conventionalCommit"`
		Issue  string `yaml:"issue"`
		SemVer string `yaml:"semVer"`
	} `yaml:"regex"`
}

var ChangelogConfig ChangelogScheme

func ChangelogRead() {
	file, err := ioutil.ReadFile("autogit.yml")
	if err != nil {
		fmt.Printf("Could not read the file due to this %s error \n", err)
	}
	ChangelogConfig = ChangelogScheme{}

	err = yaml.Unmarshal(file, &ChangelogConfig)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func ChangelogValidate() {
	if ChangelogConfig.Changelog.CommitURL == "" {
		log.Fatal("autogit.yml->Changelog.CommitUrl is empty")
	}
	if ChangelogConfig.Changelog.CommitRangeURL == "" {
		log.Fatal("autogit.yml->Changelog.CommitRangeURL is empty")
	}
	if ChangelogConfig.Changelog.IssueURL == "" {
		log.Fatal("autogit.yml->Changelog.IssueURL is empty")
	}
}

func ChangelogInit() {
	ChangelogRead()
	ChangelogValidate()
}

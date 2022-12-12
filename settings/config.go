package settings

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type ConfigScheme struct {
	Changelog struct {
		CommitURL      string `yaml:"commitUrl"`
		CommitRangeURL string `yaml:"commitRangeUrl"`
		IssueURL       string `yaml:"issueUrl"`
	} `yaml:"changelog"`
}

var Config ConfigScheme

func readConfig() {
	file, err := ioutil.ReadFile("autogit.yml")
	if err != nil {
		fmt.Printf("Could not read the file due to this %s error \n", err)
	}
	Config = ConfigScheme{}

	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func validateConfig() {
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

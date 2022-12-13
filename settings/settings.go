package settings

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

//go:embed version.txt
var Version string

var Config struct {
	Changelog ChangelogScheme `yaml:"changelog"`
	Regex     RegexScheme     `yaml:"regex"`
}

func ConfigRead() {
	file, err := ioutil.ReadFile("autogit.yml")
	if err != nil {
		fmt.Printf("Could not read the file due to this %s error \n", err)
	}

	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func init() {
	ConfigRead()
	ChangelogInit()
	RegexInit()
}

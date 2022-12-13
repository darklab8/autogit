package settings

import (
	"autogit/utils"
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"gopkg.in/yaml.v3"
)

type RegexScheme struct {
	Regex struct {
		ConventionalCommit struct {
			Header     string `yaml:"header"`
			BodyFooter string `yaml:"bodyFooter"`
		} `yaml:"conventionalCommit"`
		Issue  string `yaml:"issue"`
		SemVer string `yaml:"semVer"`
	} `yaml:"regex"`
}

var RegexConventionalCommit *regexp.Regexp
var RegexBodyFooter *regexp.Regexp
var RegexIssue *regexp.Regexp
var RegexSemVer *regexp.Regexp

var RegexConfig RegexScheme

func RegexRead() {
	file, err := ioutil.ReadFile("autogit.yml")
	if err != nil {
		fmt.Printf("Could not read the file due to this %s error \n", err)
	}
	RegexConfig = RegexScheme{}

	err = yaml.Unmarshal(file, &RegexConfig)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func RegexSetDefaults() {
	if RegexConfig.Regex.ConventionalCommit.Header == "" {
		// copied from https://gist.github.com/marcojahn/482410b728c31b221b70ea6d2c433f0c
		// type, scope, subject, the rest
		RegexConfig.Regex.ConventionalCommit.Header = `^([a-z]+)(?:\(([\w]+)\))?: (?:([ -~]+))(?:\n\n([\s -~]*)|[\n])?\z`
	}
	if RegexConfig.Regex.ConventionalCommit.BodyFooter == "" {
		// everything except : which is between 9 and ;
		RegexConfig.Regex.ConventionalCommit.BodyFooter = `(?:([ -9;-~]+))\: (?:([ -9;-~]+))`
	}
	if RegexConfig.Regex.Issue == "" {
		RegexConfig.Regex.Issue = `\#([0-9]+)`
	}
	if RegexConfig.Regex.SemVer == "" {
		// copied from https://semver.org/spec/v2.0.0.html
		RegexConfig.Regex.SemVer = `^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	}
}

func RegexCompile() {
	utils.InitRegexExpression(&RegexConventionalCommit, RegexConfig.Regex.ConventionalCommit.Header)
	utils.InitRegexExpression(&RegexBodyFooter, RegexConfig.Regex.ConventionalCommit.BodyFooter)
	utils.InitRegexExpression(&RegexIssue, RegexConfig.Regex.Issue)
	utils.InitRegexExpression(&RegexSemVer, RegexConfig.Regex.SemVer)
}

func RegexInit() {
	RegexRead()
	RegexSetDefaults()
	RegexCompile()
}

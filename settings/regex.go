package settings

import (
	"autogit/utils"
	"regexp"
)

type RegexScheme struct {
	ConventionalCommit struct {
		Header     string `yaml:"header"`
		BodyFooter string `yaml:"bodyFooter"`
	} `yaml:"conventionalCommit"`
	Issue      string `yaml:"issue"`
	SemVer     string `yaml:"semVer"`
	Prerelease string `yaml:"prerelease"`
	Validation struct {
		Scope struct {
			Lowercase string `yaml:"lowercase"`
		} `yaml:"scope"`
		Type struct {
			Lowercase string `yaml:"lowercase"`
		} `yaml:"type"`
	} `yaml:"validation"`
}

var RegexConventionalCommit *regexp.Regexp
var RegexBodyFooter *regexp.Regexp
var RegexIssue *regexp.Regexp
var RegexSemVer *regexp.Regexp
var RegexScope *regexp.Regexp
var RegexType *regexp.Regexp
var RegexPrerelease *regexp.Regexp

func RegexSetDefaults() {
	if Config.Regex.ConventionalCommit.Header == "" {
		// copied from https://gist.github.com/marcojahn/482410b728c31b221b70ea6d2c433f0c
		// type, scope, subject, the rest
		Config.Regex.ConventionalCommit.Header = `^([a-z]+)(?:\(([\w]+)\))?(\!?): (?:([ -~]+))(?:\n\n([\s -~]*)|[\n])?\z`
	}
	if Config.Regex.ConventionalCommit.BodyFooter == "" {
		// everything except : which is between 9 and ;
		Config.Regex.ConventionalCommit.BodyFooter = `(?:([ -9;-~]+))\: (?:([ -9;-~]+))`
	}
	if Config.Regex.Issue == "" {
		Config.Regex.Issue = `\#([0-9]+)`
	}
	if Config.Regex.SemVer == "" {
		// copied from https://semver.org/spec/v2.0.0.html
		Config.Regex.SemVer = `^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	}

	if Config.Regex.Validation.Scope.Lowercase == "" {
		Config.Regex.Validation.Scope.Lowercase = `^[a-z]+$`
	}

	if Config.Regex.Validation.Type.Lowercase == "" {
		Config.Regex.Validation.Type.Lowercase = `^[a-z]+$`
	}

	if Config.Regex.Prerelease == "" {
		Config.Regex.Prerelease = `(?:a\.([0-9]+))?\.?(?:b\.([0-9]+))?\.?(?:rc\.([0-9]+))?`
	}
}

func RegexCompile() {
	utils.InitRegexExpression(&RegexConventionalCommit, Config.Regex.ConventionalCommit.Header)
	utils.InitRegexExpression(&RegexBodyFooter, Config.Regex.ConventionalCommit.BodyFooter)
	utils.InitRegexExpression(&RegexIssue, Config.Regex.Issue)
	utils.InitRegexExpression(&RegexSemVer, Config.Regex.SemVer)
	utils.InitRegexExpression(&RegexScope, Config.Regex.Validation.Scope.Lowercase)
	utils.InitRegexExpression(&RegexType, Config.Regex.Validation.Type.Lowercase)
	utils.InitRegexExpression(&RegexPrerelease, Config.Regex.Prerelease)
}

func RegexInit() {
	RegexSetDefaults()
	RegexCompile()
}

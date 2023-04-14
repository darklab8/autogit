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

func RegexSetDefaults(conf *ConfigScheme) {
	if conf.Regex.ConventionalCommit.Header == "" {
		// copied from https://gist.github.com/marcojahn/482410b728c31b221b70ea6d2c433f0c
		// type, scope, subject, the rest
		conf.Regex.ConventionalCommit.Header = `^([a-z]+)(?:\(([\w]+)\))?(\!?): (?:([ -~]+))(?:\n\n([\s -~]*)|[\n])?\z`
	}
	if conf.Regex.ConventionalCommit.BodyFooter == "" {
		// everything except : which is between 9 and ;
		conf.Regex.ConventionalCommit.BodyFooter = `(?:([ -9;-~]+))\: (?:([ -9;-~]+))`
	}
	if conf.Regex.Issue == "" {
		conf.Regex.Issue = `\#([0-9]+)`
	}
	if conf.Regex.SemVer == "" {
		// copied from https://semver.org/spec/v2.0.0.html
		conf.Regex.SemVer = `^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	}

	if conf.Regex.Validation.Scope.Lowercase == "" {
		conf.Regex.Validation.Scope.Lowercase = `^[a-z]+$`
	}

	if conf.Regex.Validation.Type.Lowercase == "" {
		conf.Regex.Validation.Type.Lowercase = `^[a-z]+$`
	}

	if conf.Regex.Prerelease == "" {
		conf.Regex.Prerelease = `(?:a\.([0-9]+))?\.?(?:b\.([0-9]+))?\.?(?:rc\.([0-9]+))?`
	}
}

func RegexCompile(conf *ConfigScheme) {
	utils.InitRegexExpression(&RegexConventionalCommit, conf.Regex.ConventionalCommit.Header)
	utils.InitRegexExpression(&RegexBodyFooter, conf.Regex.ConventionalCommit.BodyFooter)
	utils.InitRegexExpression(&RegexIssue, conf.Regex.Issue)
	utils.InitRegexExpression(&RegexSemVer, conf.Regex.SemVer)
	utils.InitRegexExpression(&RegexScope, conf.Regex.Validation.Scope.Lowercase)
	utils.InitRegexExpression(&RegexType, conf.Regex.Validation.Type.Lowercase)
	utils.InitRegexExpression(&RegexPrerelease, conf.Regex.Prerelease)
}

func RegexInit(conf *ConfigScheme) {
	RegexSetDefaults(conf)
	RegexCompile(conf)
}

package settings

import (
	"autogit/settings/types"
	"autogit/settings/utils"
	"regexp"
)

type RegexScheme struct {
	ConventionalCommit struct {
		Header     types.RegexExpression `yaml:"header"`
		BodyFooter types.RegexExpression `yaml:"bodyFooter"`
	} `yaml:"conventionalCommit"`
	Issue      types.RegexExpression `yaml:"issue"`
	SemVer     types.RegexExpression `yaml:"semVer"`
	Prerelease types.RegexExpression `yaml:"prerelease"`
	Validation struct {
		Scope struct {
			Lowercase types.RegexExpression `yaml:"lowercase"`
		} `yaml:"scope"`
		Type struct {
			Lowercase types.RegexExpression `yaml:"lowercase"`
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
	RegexCompile(conf)
}

package settings

import (
	"autogit/settings/types"
	"autogit/settings/utils"
	"regexp"
)

type RegexScheme struct {
	ConventionalCommit struct {
		Header types.RegexExpression `yaml:"header"`
		Github struct {
			Enabled     bool                  `yaml:"enabled"`
			MergeHeader types.RegexExpression `yaml:"merge_header"`
		} `yaml:"github"`
		BodyFooter types.RegexExpression `yaml:"bodyFooter"`
	} `yaml:"conventionalCommit"`
	Issue            types.RegexExpression `yaml:"issue"`
	PullRequestRegex types.RegexExpression `yaml:"pull_request"`
	SemVer           types.RegexExpression `yaml:"semVer"`
	Prerelease       types.RegexExpression `yaml:"prerelease"`
	Validation       struct {
		Scope struct {
			Lowercase types.RegexExpression `yaml:"lowercase"`
		} `yaml:"scope"`
		Type struct {
			Lowercase types.RegexExpression `yaml:"lowercase"`
		} `yaml:"type"`
	} `yaml:"validation"`
}

var RegexConventionalCommit *regexp.Regexp
var RegexGithubMergeCommit *regexp.Regexp
var RegexBodyFooter *regexp.Regexp
var RegexIssue *regexp.Regexp
var RegexSemVer *regexp.Regexp
var RegexScope *regexp.Regexp
var RegexType *regexp.Regexp
var RegexPrerelease *regexp.Regexp
var RegexPullRequest *regexp.Regexp

func (conf *ConfigScheme) regexCompile() {
	utils.InitRegexExpression(&RegexConventionalCommit, conf.Regex.ConventionalCommit.Header)
	utils.InitRegexExpression(&RegexGithubMergeCommit, conf.Regex.ConventionalCommit.Github.MergeHeader)
	utils.InitRegexExpression(&RegexBodyFooter, conf.Regex.ConventionalCommit.BodyFooter)
	utils.InitRegexExpression(&RegexIssue, conf.Regex.Issue)
	utils.InitRegexExpression(&RegexSemVer, conf.Regex.SemVer)
	utils.InitRegexExpression(&RegexScope, conf.Regex.Validation.Scope.Lowercase)
	utils.InitRegexExpression(&RegexType, conf.Regex.Validation.Type.Lowercase)
	utils.InitRegexExpression(&RegexPrerelease, conf.Regex.Prerelease)
	utils.InitRegexExpression(&RegexPullRequest, conf.Regex.PullRequestRegex)
}

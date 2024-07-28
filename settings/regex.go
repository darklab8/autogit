package settings

import (
	"regexp"

	"github.com/darklab8/go-utils/utils/regexy"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type RegexScheme struct {
	ConventionalCommit struct {
		Headers    []utils_types.RegExp `yaml:"headers"`
		BodyFooter utils_types.RegExp   `yaml:"bodyFooter"`
	} `yaml:"conventionalCommit"`
	Issue            utils_types.RegExp `yaml:"issue"`
	PullRequestRegex utils_types.RegExp `yaml:"pull_request"`
	SemVer           utils_types.RegExp `yaml:"semVer"`
	Prerelease       utils_types.RegExp `yaml:"prerelease"`
	Validation       struct {
		Scope struct {
			Lowercase utils_types.RegExp `yaml:"lowercase"`
		} `yaml:"scope"`
		Type struct {
			Lowercase utils_types.RegExp `yaml:"lowercase"`
		} `yaml:"type"`
	} `yaml:"validation"`
}

var RegexConventionalCommit []*regexp.Regexp = []*regexp.Regexp{}
var RegexBodyFooter *regexp.Regexp
var RegexIssue *regexp.Regexp
var RegexSemVer *regexp.Regexp
var RegexScope *regexp.Regexp
var RegexType *regexp.Regexp
var RegexPrerelease *regexp.Regexp
var RegexPullRequest *regexp.Regexp

func (conf *ConfigScheme) regexCompile() {
	if len(RegexConventionalCommit) == 0 {
		for _, regex_expression := range conf.Regex.ConventionalCommit.Headers {
			regex := &regexp.Regexp{}
			regexy.InitRegexExpression(&regex, regex_expression)
			RegexConventionalCommit = append(RegexConventionalCommit, regex)
		}
	}

	regexy.InitRegexExpression(&RegexBodyFooter, conf.Regex.ConventionalCommit.BodyFooter)
	regexy.InitRegexExpression(&RegexIssue, conf.Regex.Issue)
	regexy.InitRegexExpression(&RegexSemVer, conf.Regex.SemVer)
	regexy.InitRegexExpression(&RegexScope, conf.Regex.Validation.Scope.Lowercase)
	regexy.InitRegexExpression(&RegexType, conf.Regex.Validation.Type.Lowercase)
	regexy.InitRegexExpression(&RegexPrerelease, conf.Regex.Prerelease)
	regexy.InitRegexExpression(&RegexPullRequest, conf.Regex.PullRequestRegex)
}

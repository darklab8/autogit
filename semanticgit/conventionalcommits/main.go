package conventionalcommits

import (
	"autogit/semanticgit/conventionalcommits/conventionalcommitstype"
	"autogit/settings"
	"autogit/settings/types"
	"fmt"
	"strings"
)

type ConventionalCommit struct {
	Original types.CommitMessage
	conventionalcommitstype.ParsedCommit
}

func (commit ConventionalCommit) StringHeader() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s", commit.Type))

	if commit.Scope != "" {
		sb.WriteString(fmt.Sprintf("(%s)", commit.Scope))
	}
	sb.WriteString(fmt.Sprintf(": %s", commit.Subject))

	return sb.String()
}

type NotParsed struct{}

func (m NotParsed) Error() string {
	return "not parsed at all"
}

func ParseCommit(msg types.CommitMessage) (*ConventionalCommit, error) {
	result := ConventionalCommit{}
	result.Original = msg
	main_match := settings.RegexConventionalCommit.FindStringSubmatch(string(msg))

	if len(main_match) == 0 {
		return nil, NotParsed{}
	}

	result.Type = main_match[1]
	result.Scope = main_match[2]
	result.Subject = main_match[4]

	if main_match[3] != "" {
		result.Exclamation = true
	}

	IssueMatch := settings.RegexIssue.FindAllStringSubmatch(main_match[5], -1)
	for _, match := range IssueMatch {
		result.Issue = append(result.Issue, match[1])
	}

	msgs := strings.Split(main_match[5], "\n\n")

	for index, msg := range msgs {
		match := settings.RegexBodyFooter.FindStringSubmatch(msg)
		if index == 0 && len(match) == 0 {
			result.Body = msg
		} else if len(match) > 0 {
			result.Footers = append(result.Footers, conventionalcommitstype.Footer{Token: match[1], Content: match[2]})
		}
	}

	return &result, nil
}

type InvalidType struct{}

func (m InvalidType) Error() string {
	return "invalid conventional commit Type"
}

func (c *ConventionalCommit) Validate() error {
	for _, type_ := range settings.GetConfig().Validation.Rules.Header.Type.Whitelist {
		if c.Type == type_ {
			return nil
		}
	}
	return InvalidType{}
}

func NewCommit(msg types.CommitMessage) (*ConventionalCommit, error) {
	commit, err := ParseCommit(msg)
	if err != nil {
		return commit, err
	}

	err = commit.Validate()
	if err != nil {
		return commit, err
	}

	return commit, nil
}

func (c ConventionalCommit) MajorChange() bool {
	return false
}

func (c ConventionalCommit) MinorChange() bool {
	return false
}

func (c ConventionalCommit) PatchChange() bool {
	return false
}

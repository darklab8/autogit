package conventionalcommits

import (
	"autogit/semanticgit/conventionalcommits/conventionalcommitstype"
	"autogit/settings"
	"autogit/settings/types"
	"autogit/settings/utils"
	"fmt"
	"strings"
)

type ConventionalCommit struct {
	Original types.CommitOriginalMsg
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

func ParseCommit(msg types.CommitOriginalMsg) (*ConventionalCommit, error) {
	result := ConventionalCommit{}
	result.Original = msg

	main_match := []string{}
	for _, header_regex := range settings.RegexConventionalCommit {
		main_match = header_regex.FindStringSubmatch(string(msg))
		if len(main_match) > 0 {
			break
		}
	}

	if len(main_match) == 0 {
		return nil, NotParsed{}
	}

	result.Type = conventionalcommitstype.Type(main_match[1])
	result.Scope = conventionalcommitstype.Scope(main_match[2])
	result.Subject = conventionalcommitstype.Subject(main_match[4])

	if main_match[3] != "" {
		result.Exclamation = true
	}

	IssueMatch := settings.RegexIssue.FindAllStringSubmatch(main_match[5], -1)
	for _, match := range IssueMatch {
		result.Issue = append(result.Issue, conventionalcommitstype.Issue(match[1]))
	}

	msgs := strings.Split(main_match[5], "\n\n")

	for index, msg := range msgs {
		match := settings.RegexBodyFooter.FindStringSubmatch(msg)
		if index == 0 && len(match) == 0 {
			result.Body = conventionalcommitstype.Body(msg)
		} else if len(match) > 0 {
			result.Footers = append(result.Footers, conventionalcommitstype.Footer{
				Token:   conventionalcommitstype.FooterToken(match[1]),
				Content: conventionalcommitstype.FooterContent(match[2]),
			})
		}
	}

	return &result, nil
}

type InvalidType struct {
	allowed_types []conventionalcommitstype.Type
}

func (m InvalidType) Error() string {
	return "invalid conventional commit Type. Allowed types:" + strings.Join(
		utils.CompL(m.allowed_types, func(x conventionalcommitstype.Type) string { return string(x) }),
		",",
	)
}

func (c *ConventionalCommit) Validate() error {
	allowed_types := settings.GetConfig().Validation.Rules.Header.Type.Allowlists.GetAllTypes()
	for _, type_ := range allowed_types {
		if c.Type == type_ {
			return nil
		}
	}
	return InvalidType{allowed_types: allowed_types}
}

func NewCommit(msg types.CommitOriginalMsg) (*ConventionalCommit, error) {
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

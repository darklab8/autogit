package conventionalcommits

import (
	"autogit/semanticgit/conventionalcommits/conventionalcommitstype"
	"autogit/settings"
	"autogit/settings/logus"

	"autogit/settings/types"
	"fmt"
	"strings"

	"github.com/darklab8/darklab_goutils/goutils/utils"
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

type NotParsedCommit struct{}

func (m NotParsedCommit) Error() string {
	return "not parsed commit at all"
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
		logus.Log.Debug("original msg:\n" + string(msg))
		return nil, NotParsedCommit{}
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

	body_with_footers := main_match[5]

	footers := settings.RegexBodyFooter.FindAllStringSubmatch(body_with_footers, -1)

	cleaned_body := body_with_footers
	for _, footer := range footers {
		// if u had better regex, it would not be needed to clean it from newlines :wink:
		token := strings.ReplaceAll(footer[1], "\n", "")
		result.Footers = append(result.Footers, conventionalcommitstype.Footer{
			Token:   conventionalcommitstype.FooterToken(token),
			Content: conventionalcommitstype.FooterContent(footer[2]),
		})
		cleaned_body = strings.Replace(cleaned_body, footer[0], "", -1)
	}

	cleaned_body_lines := strings.Split(cleaned_body, "\n")
	purrified_body_lines := []string{}
	for _, body_line := range cleaned_body_lines {
		if strings.HasPrefix(body_line, "#") {
			continue
		}

		purrified_body_lines = append(purrified_body_lines, body_line)
	}
	purrified_body := strings.Join(purrified_body_lines, "\n")
	purrified_body = strings.ReplaceAll(purrified_body, "\n\n", "\n")
	purrified_body = strings.ReplaceAll(purrified_body, "\n\n", "\n")

	result.Body = conventionalcommitstype.Body(purrified_body)

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

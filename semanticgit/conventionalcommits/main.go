package conventionalcommits

import (
	"autogit/settings"
	"fmt"
	"log"
	"strings"
)

type Footer struct {
	Token   string
	Content string
}

type ConventionalCommit struct {
	Original string

	Type        string
	Exclamation bool
	Scope       string
	Subject     string
	Body        string
	Footers     []Footer
	Hash        string
	Issue       []string
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

func (commit ConventionalCommit) StringAnnotated() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("type=%s\n", commit.Type))

	if commit.Scope != "" {
		sb.WriteString(fmt.Sprintf("scope=%s\n", commit.Scope))
	}
	sb.WriteString(fmt.Sprintf("subject=%s\n", commit.Subject))

	if commit.Body != "" {
		sb.WriteString(fmt.Sprintf("body=%s\n", commit.Body))
	}

	for index, footer := range commit.Footers {
		sb.WriteString(fmt.Sprintf("footer #%d - token: %s, content: %s\n", index, footer.Token, footer.Content))
	}

	return sb.String()
}

func Check(err error, strict bool, msgs ...string) {
	if err == nil {
		return
	}

	if !strict {
		log.Print(err, msgs)
	}

	log.Fatal(err, msgs)
}

type NotParsed struct{}

func (m NotParsed) Error() string {
	return "not parsed at all"
}

func ParseCommit(msg string) (*ConventionalCommit, error) {
	result := ConventionalCommit{}
	result.Original = msg
	main_match := settings.RegexConventionalCommit.FindStringSubmatch(msg)

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
			result.Footers = append(result.Footers, Footer{Token: match[1], Content: match[2]})
		}
	}

	return &result, nil
}

type InvalidType struct{}

func (m InvalidType) Error() string {
	return "invalid conventional commit Type"
}

func (c *ConventionalCommit) Validate() error {
	for _, type_ := range settings.Config.Validation.Rules.Header.Type.Whitelist {
		if c.Type == type_ {
			return nil
		}
	}
	return InvalidType{}
}

func NewCommit(msg string) (*ConventionalCommit, error) {
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

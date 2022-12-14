package conventionalcommits

import (
	"autogit/settings"
	"log"
	"strings"
)

type Footer struct {
	Token   string
	Content string
}

type ConventionalCommit struct {
	_Original string

	Type        string
	Exclamation bool
	Scope       string
	Subject     string
	Body        string
	Footers     []Footer
	Hash        string
	Issue       []string
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
	result._Original = msg
	main_match := settings.RegexConventionalCommit.FindStringSubmatch(msg)

	if len(main_match) == 0 {
		return nil, NotParsed{}
	}

	result.Type = main_match[1]
	result.Scope = main_match[2]
	result.Subject = main_match[3]

	IssueMatch := settings.RegexIssue.FindAllStringSubmatch(main_match[4], -1)
	for _, match := range IssueMatch {
		result.Issue = append(result.Issue, match[1])
	}

	msgs := strings.Split(main_match[4], "\n\n")

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

var validTypes = [...]string{
	"build",
	"chore",
	"ci",
	"docs",
	"feat",
	"fix",
	"perf",
	"refactor",
	"revert",
	"style",
	"test",
	"merge",
}

type InvalidType struct{}

func (m InvalidType) Error() string {
	return "invalid conventional commit Type"
}

func (c *ConventionalCommit) Validate() error {
	for _, type_ := range validTypes {
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

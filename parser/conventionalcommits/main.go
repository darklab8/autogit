package conventionalcommits

import (
	"autogit/utils"
	"log"
	"regexp"
	"strings"
)

type Footer struct {
	Token   string
	Content string
}

type ConventionalCommit struct {
	Type        string
	Exclamation bool
	Scope       string
	Subject     string
	Body        string
	Footers     []Footer
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

var conventionalRegex *regexp.Regexp
var FooterRegex *regexp.Regexp

func init() {
	// copied from https://gist.github.com/marcojahn/482410b728c31b221b70ea6d2c433f0c

	utils.InitRegexExpression(&conventionalRegex,
		// type, scope, subject, the rest
		`^([a-z]+)(?:\(([\w]+)\))?: (?:([a-zA-Z 0-9-.]+))(?:\n\n([\w\s\-\:]*))?`)

	utils.InitRegexExpression(&FooterRegex,
		`(?:([a-zA-Z 0-9-.]+))\: (?:([a-zA-Z 0-9-.]+))`)
}

type NotParsed struct{}

func (m NotParsed) Error() string {
	return "not parsed at all"
}

func ParseCommit(msg string) (*ConventionalCommit, error) {
	result := ConventionalCommit{}
	main_match := conventionalRegex.FindStringSubmatch(msg)

	if len(main_match) == 0 {
		return nil, NotParsed{}
	}

	result.Type = main_match[1]
	result.Scope = main_match[2]
	result.Subject = main_match[3]

	msgs := strings.Split(main_match[4], "\n\n")

	for index, msg := range msgs {
		match := FooterRegex.FindStringSubmatch(msg)
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

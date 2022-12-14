package validation

import (
	"autogit/semanticgit/conventionalcommits"
	"autogit/settings"
	"fmt"
)

type ErrorInvalidMaxLength struct {
	commit *conventionalcommits.ConventionalCommit
}

func (err ErrorInvalidMaxLength) Error() string {
	return fmt.Sprintf("commit.header='%s' has length greater than maxLength=%d", err.commit.StringHeader(), settings.Config.Validation.Rules.Header.MaxLength)
}

type ErrorNotFoundIssue struct {
	commit *conventionalcommits.ConventionalCommit
}

func (err ErrorNotFoundIssue) Error() string {
	return fmt.Sprintf("commit='%s' has no linked issue with regex %s", err.commit.StringHeader(), settings.Config.Regex.Issue)
}

func Validate(commit *conventionalcommits.ConventionalCommit) error {

	if len(commit.StringHeader()) > settings.Config.Validation.Rules.Header.MaxLength {
		return ErrorInvalidMaxLength{commit: commit}
	}

	if settings.Config.Validation.Rules.Issue.Present {
		IssueMatch := settings.RegexIssue.FindAllStringSubmatch(commit.Original, -1)
		if len(IssueMatch) == 0 {
			return ErrorNotFoundIssue{commit: commit}
		}

	}

	return nil
}

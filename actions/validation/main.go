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
	return fmt.Sprintf("commit.header='%s' must have length shorter than maxLength=%d, because rule is enabled", err.commit.StringHeader(), settings.Config.Validation.Rules.Header.MaxLength)
}

type ErrorNotFoundIssue struct {
	commit *conventionalcommits.ConventionalCommit
}

func (err ErrorNotFoundIssue) Error() string {
	return fmt.Sprintf("commit='%s' must have linked issue regex %s, because rule is enabled", err.commit.StringHeader(), settings.Config.Regex.Issue)
}

type ErrorCommitScopeMustBeDefined struct {
	commit *conventionalcommits.ConventionalCommit
}

func (err ErrorCommitScopeMustBeDefined) Error() string {
	return fmt.Sprintf("commit='%s' must have defined scope, because rule is enabled", err.commit.StringHeader())
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

	if settings.Config.Validation.Rules.Header.Scope.Present {
		if commit.Scope == "" {
			return ErrorCommitScopeMustBeDefined{commit: commit}
		}
	}

	return nil
}

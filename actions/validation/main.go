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
	return fmt.Sprintf("length(commit.header)=%s is greater than maxLength=%d", err.commit.StringHeader(), settings.Config.Validation.Rules.Header.MaxLength)
}

func Validate(commit *conventionalcommits.ConventionalCommit) error {

	if len(commit.StringHeader()) > settings.Config.Validation.Rules.Header.MaxLength {
		return ErrorInvalidMaxLength{commit: commit}
	}

	return nil
}

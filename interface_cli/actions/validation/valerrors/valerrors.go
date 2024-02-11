package valerrors

import (
	"fmt"

	"github.com/darklab8/autogit/semanticgit/conventionalcommits"
	"github.com/darklab8/autogit/settings"
)

type errorInvalidMaxLength struct {
	commit conventionalcommits.ConventionalCommit
	conf   settings.ConfigScheme
}

func (err errorInvalidMaxLength) Error() string {
	return fmt.Sprintf("commit.header='%s' must have length shorter than maxLength=%d, because rule is enabled", err.commit.StringHeader(), err.conf.Validation.Rules.Header.MaxLength)
}

func NewErrorInvalidMaxLength(commit conventionalcommits.ConventionalCommit,
	conf settings.ConfigScheme) errorInvalidMaxLength {
	return errorInvalidMaxLength{commit: commit, conf: conf}
}

type errorNotFoundIssue struct {
	commit conventionalcommits.ConventionalCommit
	conf   settings.ConfigScheme
}

func (err errorNotFoundIssue) Error() string {
	return fmt.Sprintf("commit='%s' must have linked issue regex %s, because rule is enabled", err.commit.StringHeader(), err.conf.Regex.Issue)
}

func NewErrorNotFoundIssue(commit conventionalcommits.ConventionalCommit,
	conf settings.ConfigScheme) errorNotFoundIssue {
	return errorNotFoundIssue{commit: commit, conf: conf}
}

type errorCommitScopeMustBeDefined struct {
	commit conventionalcommits.ConventionalCommit
}

func (err errorCommitScopeMustBeDefined) Error() string {
	return fmt.Sprintf("commit='%s' must have defined scope, because rule is enabled", err.commit.StringHeader())
}

func NewErrorCommitScopeMustBeDefined(commit conventionalcommits.ConventionalCommit) errorCommitScopeMustBeDefined {
	return errorCommitScopeMustBeDefined{commit: commit}
}

type errorCommitScopeMustBeLowercase struct {
	commit conventionalcommits.ConventionalCommit
}

func (err errorCommitScopeMustBeLowercase) Error() string {
	return fmt.Sprintf("commit='%s' must have scope in lowercase, because rule is enabled", err.commit.StringHeader())
}

func NewerrorCommitScopeMustBeLowercase(commit conventionalcommits.ConventionalCommit) errorCommitScopeMustBeLowercase {
	return errorCommitScopeMustBeLowercase{commit: commit}
}

type errorCommitTypeMustBeLowercase struct {
	commit conventionalcommits.ConventionalCommit
}

func (err errorCommitTypeMustBeLowercase) Error() string {
	return fmt.Sprintf("commit='%s' must have type in lowercase, because rule is enabled", err.commit.StringHeader())
}

func NewerrorCommitTypeMustBeLowercase(commit conventionalcommits.ConventionalCommit) errorCommitTypeMustBeLowercase {
	return errorCommitTypeMustBeLowercase{commit: commit}
}

type errorCommitScopeMustBeInAllowlist struct {
	commit conventionalcommits.ConventionalCommit
	conf   settings.ConfigScheme
}

func (err errorCommitScopeMustBeInAllowlist) Error() string {
	return fmt.Sprintf("commit='%s' must be in allowlist %v, because rule is enabled", err.commit.StringHeader(), err.conf.Validation.Rules.Header.Scope.Allowlist)
}

func NewerrorCommitScopeMustBeInAllowlist(commit conventionalcommits.ConventionalCommit, conf settings.ConfigScheme) errorCommitScopeMustBeInAllowlist {
	return errorCommitScopeMustBeInAllowlist{commit: commit, conf: conf}
}

type errorCommitSubjectMinWords struct {
	commit conventionalcommits.ConventionalCommit
	conf   settings.ConfigScheme
}

func (err errorCommitSubjectMinWords) Error() string {
	return fmt.Sprintf("commit='%s' must have in subject at least %d words, because rule is enabled", err.commit.StringHeader(), err.conf.Validation.Rules.Header.Subject.MinWords)
}

func NewerrorCommitSubjectMinWords(commit conventionalcommits.ConventionalCommit, conf settings.ConfigScheme) errorCommitSubjectMinWords {
	return errorCommitSubjectMinWords{commit: commit, conf: conf}
}

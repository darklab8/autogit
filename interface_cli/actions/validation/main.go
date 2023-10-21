package validation

import (
	"autogit/interface_cli/actions/validation/valerrors"
	"autogit/semanticgit/conventionalcommits"
	"autogit/settings"
	"strings"
)

func Validate(commit conventionalcommits.ConventionalCommit, conf settings.ConfigScheme) error {

	if len(commit.StringHeader()) > conf.Validation.Rules.Header.MaxLength {
		return valerrors.NewErrorInvalidMaxLength(commit, conf)
	}

	if conf.Validation.Rules.Issue.Present {
		IssueMatch := settings.RegexIssue.FindAllStringSubmatch(commit.Original.ToString(), -1)
		if len(IssueMatch) == 0 {
			return valerrors.NewErrorNotFoundIssue(commit, conf)
		}
	}

	if conf.Validation.Rules.Header.Scope.Present {
		if commit.Scope == "" {
			return valerrors.NewErrorCommitScopeMustBeDefined(commit)
		}

		if conf.Validation.Rules.Header.Scope.Lowercase {
			if !settings.RegexScope.MatchString(commit.Scope) {
				return valerrors.NewerrorCommitScopeMustBeLowercase(commit)
			}
		}
	}

	if conf.Validation.Rules.Header.Type.Lowercase {
		if !settings.RegexType.MatchString(commit.Type) {
			return valerrors.NewerrorCommitTypeMustBeLowercase(commit)
		}
	}

	if len(conf.Validation.Rules.Header.Scope.Whitelist) > 0 && commit.Scope != "" {
		matchFound := false
		for _, allowed_scope := range conf.Validation.Rules.Header.Scope.Whitelist {
			if commit.Scope == allowed_scope {
				matchFound = true
			}
		}

		if !matchFound {
			return valerrors.NewerrorCommitScopeMustBeInWhitelist(commit, conf)
		}
	}

	words := strings.Split(commit.Subject, " ")
	if len(words) < conf.Validation.Rules.Header.Subject.MinWords {
		return valerrors.NewerrorCommitSubjectMinWords(commit, conf)
	}

	return nil
}

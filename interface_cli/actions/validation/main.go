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

	for _, commit_type := range conf.Validation.Rules.Header.Scope.EnforcedForTypes {
		if commit.Type != commit_type {
			continue
		}

		if commit.Scope == "" {
			return valerrors.NewErrorCommitScopeMustBeDefined(commit)
		}
	}

	if commit.Scope != "" {
		if conf.Validation.Rules.Header.Scope.Lowercase {
			if !settings.RegexScope.MatchString(string(commit.Scope)) {
				return valerrors.NewerrorCommitScopeMustBeLowercase(commit)
			}
		}
	}

	if conf.Validation.Rules.Header.Type.Lowercase {
		if !settings.RegexType.MatchString(string(commit.Type)) {
			return valerrors.NewerrorCommitTypeMustBeLowercase(commit)
		}
	}

	if len(conf.Validation.Rules.Header.Scope.Allowlist) > 0 && commit.Scope != "" {
		matchFound := false
		for _, allowed_scope := range conf.Validation.Rules.Header.Scope.Allowlist {
			if commit.Scope == allowed_scope {
				matchFound = true
			}
		}

		if !matchFound {
			return valerrors.NewerrorCommitScopeMustBeInAllowlist(commit, conf)
		}
	}

	words := strings.Split(string(commit.Subject), " ")
	if len(words) < conf.Validation.Rules.Header.Subject.MinWords {
		return valerrors.NewerrorCommitSubjectMinWords(commit, conf)
	}

	return nil
}

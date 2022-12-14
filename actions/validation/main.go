package validation

import "autogit/semanticgit/conventionalcommits"

type Validator struct {
	Commits []*conventionalcommits.ConventionalCommit
}

func (v Validator) Run() {

}

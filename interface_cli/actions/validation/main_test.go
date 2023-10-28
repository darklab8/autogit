package validation

import (
	"autogit/settings/testutils"
	_ "autogit/settings/testutils/autouse"

	"autogit/semanticgit/conventionalcommits"
	"autogit/semanticgit/conventionalcommits/conventionalcommitstype"
	"autogit/settings"
	"autogit/settings/logus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxLengthHeaderErrorNo(t *testing.T) {
	conf := &settings.ConfigScheme{}
	conf.Validation.Rules.Header.Type.Allowlists.SemVerMinorIncreasers = []conventionalcommitstype.Type{"feat"}
	conf.Validation.Rules.Header.MaxLength = 72
	commit, err := conventionalcommits.NewCommit("feat: abc")
	logus.CheckFatal(err, "failed creating commit")
	err = Validate(*commit, *conf)
	testutils.Equal(t, nil, err)
}

func TestMaxLengthHeaderErrorYes(t *testing.T) {
	conf := &settings.ConfigScheme{}
	conf.Validation.Rules.Header.Type.Allowlists.SemVerMinorIncreasers = []conventionalcommitstype.Type{"feat"}
	conf.Validation.Rules.Header.MaxLength = 72
	commit, err := conventionalcommits.NewCommit("feat: writing long on purpose commit, which should be way beyond 72 characters")
	logus.CheckFatal(err, "failed creating commit")
	err = Validate(*commit, *conf)
	assert.NotEqual(t, nil, err)
}

func TestScopeEnforcing(t *testing.T) {
	conf := &settings.ConfigScheme{}
	conf.Validation.Rules.Header.MaxLength = 72

	conf.Validation.Rules.Header.Scope.EnforcedForTypes = []conventionalcommitstype.Type{}
	commit, err := conventionalcommits.NewCommit("feat: some common commit")
	logus.CheckFatal(err, "failed creating commit")
	err = Validate(*commit, *conf)
	assert.Nil(t, err)

	conf.Validation.Rules.Header.Scope.EnforcedForTypes = []conventionalcommitstype.Type{"feat"}
	commit, err = conventionalcommits.NewCommit("feat: some common commit")
	logus.CheckFatal(err, "failed creating commit")
	err = Validate(*commit, *conf)
	assert.Error(t, err)

	commit, err = conventionalcommits.NewCommit("feat(api): some common commit")
	logus.CheckFatal(err, "failed creating commit")
	err = Validate(*commit, *conf)
	assert.Nil(t, err)

	conf.Validation.Rules.Header.Scope.Allowlist = []conventionalcommitstype.Scope{"smth"}

	commit, err = conventionalcommits.NewCommit("feat(api): some common commit")
	logus.CheckFatal(err, "failed creating commit")
	err = Validate(*commit, *conf)
	assert.Error(t, err)

	commit, err = conventionalcommits.NewCommit("feat(smth): some common commit")
	logus.CheckFatal(err, "failed creating commit")
	err = Validate(*commit, *conf)
	assert.Nil(t, err)
}

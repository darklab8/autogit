package validation

import (
	_ "autogit/settings/testutils/autouse"

	"autogit/semanticgit/conventionalcommits"
	"autogit/settings"
	"autogit/settings/logus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxLengthHeaderErrorNo(t *testing.T) {
	conf := &settings.ConfigScheme{}
	conf.Validation.Rules.Header.Type.Whitelist = []string{"feat"}
	conf.Validation.Rules.Header.MaxLength = 72
	commit, err := conventionalcommits.NewCommit("feat: abc")
	logus.CheckFatal(err, "failed creating commit")
	err = Validate(commit, conf)
	assert.Equal(t, nil, err)
}

func TestMaxLengthHeaderErrorYes(t *testing.T) {
	conf := &settings.ConfigScheme{}
	conf.Validation.Rules.Header.Type.Whitelist = []string{"feat"}
	conf.Validation.Rules.Header.MaxLength = 72
	commit, err := conventionalcommits.NewCommit("feat: writing long on purpose commit, which should be way beyond 72 characters")
	logus.CheckFatal(err, "failed creating commit")
	err = Validate(commit, conf)
	assert.NotEqual(t, nil, err)
}

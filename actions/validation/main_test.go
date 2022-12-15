package validation

import (
	"autogit/semanticgit/conventionalcommits"
	"autogit/settings"
	"autogit/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxLengthHeaderErrorNo(t *testing.T) {
	settings.Config = settings.ConfigScheme{}
	settings.Config.Validation.Rules.Header.Type.Whitelist = []string{"feat"}
	settings.Config.Validation.Rules.Header.MaxLength = 72
	commit, err := conventionalcommits.NewCommit("feat: abc")
	utils.CheckPanic(err)
	err = Validate(commit)
	assert.Equal(t, nil, err)
}

func TestMaxLengthHeaderErrorYes(t *testing.T) {
	settings.Config = settings.ConfigScheme{}
	settings.Config.Validation.Rules.Header.Type.Whitelist = []string{"feat"}
	settings.Config.Validation.Rules.Header.MaxLength = 72
	commit, err := conventionalcommits.NewCommit("feat: writing long on purpose commit, which should be way beyond 72 characters")
	utils.CheckPanic(err)
	err = Validate(commit)
	assert.NotEqual(t, nil, err)
}

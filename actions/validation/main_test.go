package validation

import (
	"autogit/semanticgit/conventionalcommits"
	"autogit/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxLengthHeaderErrorNo(t *testing.T) {
	commit, err := conventionalcommits.NewCommit("feat: abc")
	utils.CheckPanic(err)
	err = Validate(commit)
	assert.Equal(t, nil, err)
}

func TestMaxLengthHeaderErrorYes(t *testing.T) {
	commit, err := conventionalcommits.NewCommit("feat: writing long on purpose commit, which should be way beyond 72 characters")
	utils.CheckPanic(err)
	err = Validate(commit)
	assert.NotEqual(t, nil, err)
}

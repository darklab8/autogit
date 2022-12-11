package conventionalcommits

import (
	"autogit/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	_, err := ParseCommit("feat: abc")
	utils.CheckPanic(err)
}

func TestParse2(t *testing.T) {
	_, err := ParseCommit("fix: dsfsdf")
	utils.CheckPanic(err)
}

func TestParse3(t *testing.T) {
	_, err := ParseCommit("abcd abc")
	assert.True(t, err != nil)
}

func TestParse4(t *testing.T) {
	commit, err := ParseCommit(`feat(api): my subject

body message

footer-key: tralala`)
	utils.CheckPanic(err)

	assert.Equal(t, "feat", commit.Type)
	assert.Equal(t, "api", commit.Scope)
	assert.Equal(t, "my subject", commit.Subject)
	assert.Equal(t, "body message", commit.Body)
	assert.Equal(t, "footer-key", commit.Footers[0].Token)
	assert.Equal(t, "tralala", commit.Footers[0].Content)
}

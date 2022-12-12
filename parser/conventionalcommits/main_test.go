package conventionalcommits

import (
	"autogit/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	_, err := NewCommit("feat: abc")
	utils.CheckPanic(err)
}

func TestParse2(t *testing.T) {
	_, err := NewCommit("fix: dsfsdf")
	utils.CheckPanic(err)
}

func TestParse3(t *testing.T) {
	_, err := NewCommit("abcd abc")
	assert.True(t, err != nil)
}

func TestParse4(t *testing.T) {
	commit, err := NewCommit(`feat(api): my subject

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

func TestParse5(t *testing.T) {
	_, err := NewCommit("feat: hook commit msg")
	utils.CheckPanic(err)
}

func TestNotAllowedNewLineAtSecondLine(t *testing.T) {
	_, err := NewCommit(`feat(api): my subject
body message

footer-key: tralala`)
	assert.True(t, err != nil)
}

func TestParse6(t *testing.T) {
	_, err := NewCommit(`refactor: autogit version into about`)
	utils.CheckPanic(err)
}

func TestParse7(t *testing.T) {
	commit, err := NewCommit(`feat: rendered changelog for task #1, #2, #3`)
	utils.CheckPanic(err)
	assert.Equal(t, "1", commit.Issue)
}

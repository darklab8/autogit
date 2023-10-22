package conventionalcommits_test

import (
	"autogit/semanticgit/conventionalcommits"
	"autogit/settings/logus"
	_ "autogit/settings/testutils/autouse"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	_, err := conventionalcommits.NewCommit("feat: abc")
	logus.CheckFatal(err, "failled creating commit")
}

func TestParse2(t *testing.T) {
	_, err := conventionalcommits.NewCommit("fix: dsfsdf")
	logus.CheckFatal(err, "failed creating commit")
}

func TestParse3(t *testing.T) {
	_, err := conventionalcommits.NewCommit("abcd abc")
	assert.True(t, err != nil)
}

func TestParse4(t *testing.T) {
	commit, err := conventionalcommits.NewCommit(`feat(api): my subject

body message

footer-key: tralala`)
	logus.CheckFatal(err, "failed creating commit")

	assert.Equal(t, "feat", commit.Type)
	assert.Equal(t, "api", commit.Scope)
	assert.Equal(t, "my subject", commit.Subject)
	assert.Equal(t, "body message", commit.Body)
	assert.Equal(t, "footer-key", commit.Footers[0].Token)
	assert.Equal(t, "tralala", commit.Footers[0].Content)
}

func TestParse5(t *testing.T) {
	_, err := conventionalcommits.NewCommit("feat: hook commit msg")
	logus.CheckFatal(err, "failed creating commit")
}

func TestNotAllowedNewLineAtSecondLine(t *testing.T) {
	_, err := conventionalcommits.NewCommit(`feat(api): my subject
body message

footer-key: tralala`)
	assert.True(t, err != nil)
}

func TestParse6(t *testing.T) {
	_, err := conventionalcommits.NewCommit(`refactor: autogit semver into about`)
	logus.CheckFatal(err, "failed creating commit")
}

func TestParse7(t *testing.T) {
	commit, err := conventionalcommits.NewCommit(`feat: rendered changelog for task

#1, #2, #3`)
	logus.CheckFatal(err, "failed creating commits")
	assert.Equal(t, "1", commit.Issue[0])
	assert.Equal(t, "2", commit.Issue[1])
	assert.Equal(t, "3", commit.Issue[2])
}

package conventionalcommits_test

import (
	"github.com/darklab8/autogit/semanticgit/conventionalcommits"
	"github.com/darklab8/autogit/semanticgit/conventionalcommits/conventionalcommitstype"
	"github.com/darklab8/autogit/settings/logus"
	"github.com/darklab8/autogit/settings/testutils"
	_ "github.com/darklab8/autogit/settings/testutils/autouse"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	_, err := conventionalcommits.NewCommit("feat: abc")
	logus.Log.CheckFatal(err, "failled creating commit")
}

func TestParse2(t *testing.T) {
	_, err := conventionalcommits.NewCommit("fix: dsfsdf")
	logus.Log.CheckFatal(err, "failed creating commit")
}

func TestParse3(t *testing.T) {
	_, err := conventionalcommits.NewCommit("abcd abc")
	assert.True(t, err != nil)
}

func TestParse4(t *testing.T) {
	commit, err := conventionalcommits.NewCommit(`feat(api): my subject

body message

footer-key: tralala`)
	logus.Log.CheckFatal(err, "failed creating commit")

	testutils.Equal(t, "feat", commit.Type)
	testutils.Equal(t, conventionalcommitstype.Scope("api"), commit.Scope)
	testutils.Equal(t, "my subject", commit.Subject)
	testutils.Equal(t, "body message\n", commit.Body)
	testutils.Equal(t, "footer-key", commit.Footers[0].Token)
	testutils.Equal(t, "tralala", commit.Footers[0].Content)
}

func TestParse5(t *testing.T) {
	_, err := conventionalcommits.NewCommit("feat: hook commit msg")
	logus.Log.CheckFatal(err, "failed creating commit")
}

func TestNotAllowedNewLineAtSecondLine(t *testing.T) {
	_, err := conventionalcommits.NewCommit(`feat(api): my subject
body message

footer-key: tralala`)
	assert.True(t, err != nil)
}

func TestParse6(t *testing.T) {
	_, err := conventionalcommits.NewCommit(`refactor: autogit semver into about`)
	logus.Log.CheckFatal(err, "failed creating commit")
}

func TestParse7(t *testing.T) {
	commit, err := conventionalcommits.NewCommit(`feat: rendered changelog for task

i-#1, i-#2, i-#3`)
	logus.Log.CheckFatal(err, "failed creating commits")
	testutils.Equal(t, "1", commit.Issue[0])
	testutils.Equal(t, "2", commit.Issue[1])
	testutils.Equal(t, "3", commit.Issue[2])
}

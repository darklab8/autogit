package changelog

import (
	"autogit/settings"
	"autogit/settings/testutils"
	_ "autogit/settings/testutils/autouse"
	"autogit/settings/types"

	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"autogit/semanticgit/semver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitGood(t *testing.T) {
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := semanticgit.NewSemanticRepo(gitInMemory)

	gitInMemory.TestCommit("feat: test")

	testutils.EqualTag(t, "v0.0.0", gitSemantic.GetCurrentVersion().ToString())

	gitInMemory.TestCreateTag("v0.0.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: test2")

	testutils.EqualTag(t, "v0.0.1", gitSemantic.GetCurrentVersion().ToString())

	rendered := NewChangelog(gitSemantic, semver.OptionsSemVer{}, settings.GetConfig().Changelog, types.TagName("")).Render()
	assert.Contains(t, rendered, "v0.1.0")

	// historing render
	rendered = NewChangelog(gitSemantic, semver.OptionsSemVer{}, settings.GetConfig().Changelog, types.TagName("v0.0.1")).Render()
	assert.Contains(t, rendered, "v0.0.1")
}

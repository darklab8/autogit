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

func FixtureGitSemantic(t *testing.T) (*semanticgit.SemanticGit, settings.ConfigScheme) {
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := semanticgit.NewSemanticRepo(gitInMemory)

	gitInMemory.TestCommit("feat: test")

	testutils.Equal(t, "v0.0.0", gitSemantic.GetCurrentVersion().ToString())

	gitInMemory.TestCreateTag("v0.0.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: test2")

	testutils.Equal(t, "v0.0.1", gitSemantic.GetCurrentVersion().ToString())
	return gitSemantic, settings.GetConfig()
}

func TestPrepare(t *testing.T) {
	gitSemantic, config := FixtureGitSemantic(t)

	// Just for debug
	rendered := NewChangelog(gitSemantic, semver.OptionsSemVer{}, config, types.TagName(""))
	_ = rendered
	rendered = NewChangelog(gitSemantic, semver.OptionsSemVer{}, config, types.TagName("v0.0.1"))
	_ = rendered
}

func TestRender(t *testing.T) {
	gitSemantic, config := FixtureGitSemantic(t)

	rendered := NewChangelog(gitSemantic, semver.OptionsSemVer{}, config, types.TagName("")).Render()
	assert.Contains(t, rendered, "v0.1.0")

	// historing render
	rendered = NewChangelog(gitSemantic, semver.OptionsSemVer{}, config, types.TagName("v0.0.1")).Render()
	assert.Contains(t, rendered, "v0.0.1")
}

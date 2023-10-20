package changelog

import (
	"autogit/settings/testutils"
	_ "autogit/testsautouse"

	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"autogit/semanticgit/semver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitGood(t *testing.T) {
	gitInMemory := (&git.Repository{}).TestNewRepo()
	gitSemantic := (&semanticgit.SemanticGit{}).NewRepo(gitInMemory)

	gitInMemory.TestCommit("feat: test")

	testutils.EqualTag(t, "v0.0.0", gitSemantic.GetCurrentVersion().ToString())

	gitInMemory.TestCreateTag("v0.0.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: test2")

	testutils.EqualTag(t, "v0.0.1", gitSemantic.GetCurrentVersion().ToString())

	rendered := ChangelogData{}.New(gitSemantic, semver.OptionsSemVer{}).Render()
	assert.Contains(t, rendered, "v0.1.0")

	// historing render
	rendered = ChangelogData{Tag: "v0.0.1"}.New(gitSemantic, semver.OptionsSemVer{}).Render()
	assert.Contains(t, rendered, "v0.0.1")
}

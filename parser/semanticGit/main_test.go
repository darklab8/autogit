package semanticgit

import (
	"autogit/git"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitGood(t *testing.T) {
	gitInMemory := (&git.Repository{}).NewRepoTest()
	gitSemantic := (&SemanticGit{}).NewRepo(gitInMemory)

	gitInMemory.TestCommit("feat: test")

	assert.Equal(t, "v0.0.0", gitSemantic.GetCurrentVersion().ToString())

	gitInMemory.TestCreateTag("v0.0.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: test2")

	assert.Equal(t, "v0.0.1", gitSemantic.GetCurrentVersion().ToString())
}

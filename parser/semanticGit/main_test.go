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

	vers := gitSemantic.GetCurrentVersion()
	assert.Equal(t, 0, vers.Major)
	assert.Equal(t, 0, vers.Minor)
	assert.Equal(t, 0, vers.Patch)
}

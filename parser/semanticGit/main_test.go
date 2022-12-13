package semanticgit

import (
	"autogit/git"
	"testing"

	"github.com/stretchr/testify/assert"
)

// For debug on workdir
// func TestIngeration(t *testing.T) {
// 	git := (&git.Repository{}).NewRepoIntegration()
// 	gitSemantic := (&SemanticGit{}).NewRepo(git)
// 	vers := gitSemantic.CalculateNextVersion(gitSemantic.GetCurrentVersion()).ToString()
// 	assert.Equal(t, "v0.2.0", vers)
// }

func TestGitGood(t *testing.T) {
	gitInMemory := (&git.Repository{}).NewRepoTest()
	gitSemantic := (&SemanticGit{}).NewRepo(gitInMemory)

	gitInMemory.TestCommit("feat: test")

	assert.Equal(t, "v0.0.0", gitSemantic.GetCurrentVersion().ToString())

	gitInMemory.TestCreateTag("v0.0.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: test2")

	assert.Equal(t, "v0.0.1", gitSemantic.GetCurrentVersion().ToString())

	assert.Equal(t, "v0.1.0", gitSemantic.CalculateNextVersion(gitSemantic.GetCurrentVersion()).ToString())

	gitInMemory.TestCreateTag("v0.1.0", gitInMemory.TestCommit("fix: thing"))

	gitInMemory.TestCommit("fix: test2")
	assert.Equal(t, "v0.1.1", gitSemantic.CalculateNextVersion(gitSemantic.GetCurrentVersion()).ToString())

	gitInMemory.TestCommit("feat: test2")

	assert.Equal(t, "v0.2.0", gitSemantic.CalculateNextVersion(gitSemantic.GetCurrentVersion()).ToString())
}

package semanticgit

import (
	_ "autogit/testsautouse"

	"autogit/semanticgit/git"
	"autogit/semanticgit/semver"
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

func TestCurrentNextRegularVersion(t *testing.T) {
	gitInMemory := (&git.Repository{}).TestNewRepo()
	gitSemantic := (&SemanticGit{}).NewRepo(gitInMemory)
	gitInMemory.TestCommit("feat: init")

	assert.Equal(t, "v0.0.0", gitSemantic.GetCurrentVersion().ToString())

	gitInMemory.TestCreateTag("v0.0.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: test2")

	assert.Equal(t, "v0.0.1", gitSemantic.GetCurrentVersion().ToString())

	assert.Equal(t, "v0.1.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())

	gitInMemory.TestCreateTag("v0.1.0", gitInMemory.TestCommit("fix: thing"))

	gitInMemory.TestCommit("fix: test2")
	assert.Equal(t, "v0.1.1", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())

	gitInMemory.TestCommit("feat: test2")

	assert.Equal(t, "v0.2.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())

	// Semantic version should be same if no new comments
	gitInMemory.TestCreateTag("v0.2.0", gitInMemory.TestCommit("feat: new thing"))
	assert.Equal(t, "v0.2.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())
}

func TestGetChangelogs(t *testing.T) {
	gitInMemory := (&git.Repository{}).TestNewRepo()
	gitSemantic := (&SemanticGit{}).NewRepo(gitInMemory)
	gitInMemory.TestCommit("feat: init")

	gitInMemory.TestCommit("feat: test3")
	gitInMemory.TestCommit("feat: test5")
	gitInMemory.TestCreateTag("v0.0.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat(api): test")
	gitInMemory.TestCreateTag("v0.0.2", gitInMemory.TestCommit("feat(api): test2"))
	gitInMemory.TestCommit("fix: test1")
	gitInMemory.TestCommit("fix: test2")
	gitInMemory.TestCommit("fix: test3")

	logs1 := gitSemantic.GetChangelogByTag("", true)
	assert.Len(t, logs1, 3)

	logs2 := gitSemantic.GetChangelogByTag("v0.0.2", true)
	assert.Len(t, logs2, 2)

	logs3 := gitSemantic.GetChangelogByTag("v0.0.1", true)
	assert.Len(t, logs3, 4)
}

func TestTestPrereleaseVersions(t *testing.T) {
	gitInMemory := (&git.Repository{}).TestNewRepo()
	gitSemantic := (&SemanticGit{}).NewRepo(gitInMemory)

	gitInMemory.TestCommit("feat: init")
	assert.Equal(t, "v0.1.0-a.1", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true}).ToString())

	gitInMemory.TestCreateTag("v0.1.0-a.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: thing")
	assert.Equal(t, "v0.1.0-a.2", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true}).ToString())

	gitInMemory.TestCommit("feat: test5")
	gitInMemory.TestCreateTag("v0.1.0", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: thing")
	assert.Equal(t, "v0.2.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())
	assert.Equal(t, "v0.2.0-a.1.b.1", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0-a.1", gitInMemory.TestCommit("fix: thing"))
	assert.Equal(t, "v0.2.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())
	assert.Equal(t, "v0.2.0-a.1.b.1", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0-a.1.b.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("fix: thing")
	assert.Equal(t, "v0.2.0-a.2.b.2", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0-a.1.b.2", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("fix: thing")
	assert.Equal(t, "v0.2.0-a.2.b.3", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0-rc.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("fix: thing")
	assert.Equal(t, "v0.2.0-rc.2", gitSemantic.GetNextVersion(semver.OptionsSemVer{Rc: true}).ToString())
	assert.Equal(t, "v0.2.0-a.2.b.3", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0", gitInMemory.TestCommit("feat: thing"))
	gitInMemory.TestCommit("feat: thing")
	assert.Equal(t, "v0.3.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())
	assert.Equal(t, "v0.3.0-rc.1", gitSemantic.GetNextVersion(semver.OptionsSemVer{Rc: true}).ToString())
	gitInMemory.TestCreateTag("v0.3.0-rc.1", gitInMemory.TestCommit("fix: thing"))
	assert.Equal(t, "v0.3.0-rc.1", gitSemantic.GetNextVersion(semver.OptionsSemVer{Rc: true}).ToString())
	gitInMemory.TestCommit("fix: thing")
	assert.Equal(t, "v0.3.0-rc.2", gitSemantic.GetNextVersion(semver.OptionsSemVer{Rc: true}).ToString())

}

func TestBreakingChanges(t *testing.T) {
	gitInMemory := (&git.Repository{}).TestNewRepo()
	gitSemantic := (&SemanticGit{}).NewRepo(gitInMemory)
	gitInMemory.TestCommit("feat: thing")
	gitInMemory.TestCommit("feat!: break")
	assert.Equal(t, "v0.1.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())
	assert.Equal(t, "v1.0.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{Publish: true}).ToString())

	gitInMemory.TestCreateTag("v1.0.0", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat!: break")
	assert.Equal(t, "v2.0.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())

	gitInMemory.TestCreateTag("v2.0.0", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat!: break")
	assert.Equal(t, "v3.0.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())
}

func TestBuildData(t *testing.T) {
	gitInMemory := (&git.Repository{}).TestNewRepo()
	gitSemantic := (&SemanticGit{}).NewRepo(gitInMemory)
	gitInMemory.TestCommit("feat: thing")
	assert.Equal(t, "v0.1.0+123", gitSemantic.GetNextVersion(semver.OptionsSemVer{Build: "123"}).ToString())
}

func TestBug(t *testing.T) {
	gitInMemory := (&git.Repository{}).TestNewRepo()
	gitSemantic := (&SemanticGit{}).NewRepo(gitInMemory)
	gitInMemory.TestCommit("feat: thing")
	gitInMemory.TestCreateTag("v0.2.0", gitInMemory.TestCommit("feat: thing"))
	gitInMemory.TestCreateTag("v0.3.0", gitInMemory.TestCommit("feat: thing"))
	gitInMemory.TestCreateTag("v0.3.0-rc.1", gitInMemory.TestCommit("feat: thing"))
	gitInMemory.TestCreateTag("v0.3.0-rc.2", gitInMemory.TestCommit("feat: thing"))
	assert.Equal(t, "v0.4.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())
}

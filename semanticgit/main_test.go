package semanticgit

import (
	_ "github.com/darklab8/autogit/settings/testutils/autouse"

	"testing"

	"github.com/darklab8/autogit/semanticgit/git"
	"github.com/darklab8/autogit/semanticgit/semver/semvertype"
	"github.com/darklab8/autogit/settings/testutils"

	"github.com/stretchr/testify/assert"
)

// For debug on workdir
// func TestIngeration(t *testing.T) {
// 	git := git.NewRepoInWorkDir(git.SshPath(settings.GetConfig().Git.SSHPath))
// 	gitSemantic := NewSemanticRepo(git)
// 	vers := gitSemantic.CalculateNextVersion(gitSemantic.GetCurrentVersion()).ToString()
// 	testutils.EqualTag(t, "v0.2.0", vers)
// }

func TestCurrentNextRegularVersion(t *testing.T) {
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := NewSemanticRepo(gitInMemory)
	gitInMemory.TestCommit("feat: init")

	testutils.Equal(t, "v0.0.0", gitSemantic.GetCurrentVersion().ToString())

	gitInMemory.TestCreateTag("v0.0.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: test2")

	testutils.Equal(t, "v0.0.1", gitSemantic.GetCurrentVersion().ToString())

	testutils.Equal(t, "v0.1.0", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{}).ToString())

	gitInMemory.TestCreateTag("v0.1.0", gitInMemory.TestCommit("fix: thing"))

	gitInMemory.TestCommit("fix: test2")
	testutils.Equal(t, "v0.1.1", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{}).ToString())

	gitInMemory.TestCommit("feat: test2")

	testutils.Equal(t, "v0.2.0", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{}).ToString())

	// Semantic version should be same if no new comments
	gitInMemory.TestCreateTag("v0.2.0", gitInMemory.TestCommit("feat: new thing"))
	testutils.Equal(t, "v0.2.0", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{}).ToString())
}

func TestGetChangelogs(t *testing.T) {
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := NewSemanticRepo(gitInMemory)
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
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := NewSemanticRepo(gitInMemory)

	gitInMemory.TestCommit("feat: init")
	testutils.Equal(t, "v0.1.0-a.1", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Alpha: true}).ToString())

	gitInMemory.TestCreateTag("v0.1.0-a.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: thing")
	testutils.Equal(t, "v0.1.0-a.2", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Alpha: true}).ToString())

	gitInMemory.TestCommit("feat: test5")
	gitInMemory.TestCreateTag("v0.1.0", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: thing")
	testutils.Equal(t, "v0.2.0", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{}).ToString())
	testutils.Equal(t, "v0.2.0-a.1.b.1", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0-a.1", gitInMemory.TestCommit("fix: thing"))
	testutils.Equal(t, "v0.2.0", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{}).ToString())
	testutils.Equal(t, "v0.2.0-a.1.b.1", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0-a.1.b.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("fix: thing")
	testutils.Equal(t, "v0.2.0-a.2.b.2", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0-a.1.b.2", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("fix: thing")
	testutils.Equal(t, "v0.2.0-a.2.b.3", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0-rc.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("fix: thing")
	testutils.Equal(t, "v0.2.0-rc.2", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Rc: true}).ToString())
	testutils.Equal(t, "v0.2.0-a.2.b.3", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0", gitInMemory.TestCommit("feat: thing"))
	gitInMemory.TestCommit("feat: thing")
	testutils.Equal(t, "v0.3.0", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{}).ToString())
	testutils.Equal(t, "v0.3.0-rc.1", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Rc: true}).ToString())
	gitInMemory.TestCreateTag("v0.3.0-rc.1", gitInMemory.TestCommit("fix: thing"))
	testutils.Equal(t, "v0.3.0-rc.1", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Rc: true}).ToString())
	gitInMemory.TestCommit("fix: thing")
	testutils.Equal(t, "v0.3.0-rc.2", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Rc: true}).ToString())

}

func TestBreakingChanges(t *testing.T) {
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := NewSemanticRepo(gitInMemory)
	gitInMemory.TestCommit("feat: thing")
	gitInMemory.TestCommit("feat!: break")
	testutils.Equal(t, "v0.1.0", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{}).ToString())
	testutils.Equal(t, "v1.0.0", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Publish: true}).ToString())

	gitInMemory.TestCreateTag("v1.0.0", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat!: break")
	testutils.Equal(t, "v2.0.0", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{}).ToString())

	gitInMemory.TestCreateTag("v2.0.0", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat!: break")
	testutils.Equal(t, "v3.0.0", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{}).ToString())
}

func TestBuildData(t *testing.T) {
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := NewSemanticRepo(gitInMemory)
	gitInMemory.TestCommit("feat: thing")
	testutils.Equal(t, "v0.1.0+123", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{Build: "123"}).ToString())
}

func TestBug(t *testing.T) {
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := NewSemanticRepo(gitInMemory)
	gitInMemory.TestCommit("feat: thing")
	gitInMemory.TestCreateTag("v0.2.0", gitInMemory.TestCommit("feat: thing"))
	gitInMemory.TestCreateTag("v0.3.0", gitInMemory.TestCommit("feat: thing"))
	gitInMemory.TestCreateTag("v0.3.0-rc.1", gitInMemory.TestCommit("feat: thing"))
	gitInMemory.TestCreateTag("v0.3.0-rc.2", gitInMemory.TestCommit("feat: thing"))
	testutils.Equal(t, "v0.4.0", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{}).ToString())
}

func TestParseWithoutPatchVers(t *testing.T) {
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := NewSemanticRepo(gitInMemory)
	gitInMemory.TestCommit("feat: thing")
	gitInMemory.TestCreateTag("v0.2", gitInMemory.TestCommit("feat: thing"))
	gitInMemory.TestCreateTag("v0.3", gitInMemory.TestCommit("feat: thing"))

	testutils.Equal(t, "v0.3.0", gitSemantic.GetNextVersion(semvertype.OptionsSemVer{}).ToString())
}

func TestParseGHMergingCommits(t *testing.T) {
	// And commit must maintain GH link to PR like #19
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := NewSemanticRepo(gitInMemory)
	gitInMemory.TestCommit("Merge pull request #19 from Company/feat/allow_only_certain_view_routes\n\nfeat: add specifiying which view routes are allowed")

	logs1 := gitSemantic.GetChangelogByTag("", true)
	assert.Len(t, logs1, 1)

	gitInMemory.TestCommit("feat do not allow to parse this stuff")
	logs1 = gitSemantic.GetChangelogByTag("", true)
	assert.Len(t, logs1, 1)
}

func TestParseBreakingChange1(t *testing.T) {
	// And commit must maintain GH link to PR like #19
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := NewSemanticRepo(gitInMemory)
	gitInMemory.TestCommit(`feat: smth is breaking

BREAKING CHANGE: first thing
	multiline first continued
BREAKING CHANGE: second thing
BREAKING CHANGE: third thing

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# Date:      Sun Oct 29 02:38:37 2023 +0100
#
# On branch master
# Changes to be committed:
#       modified:   README.md
#
# Changes not staged for commit:
#       deleted:    autogit.yml
#`)

	logs1 := gitSemantic.GetChangelogByTag("", true)
	assert.Len(t, logs1, 1)

	assert.Len(t, logs1[0].Footers, 3)

	testutils.Equal(t, "first thing\n\tmultiline first continued", logs1[0].Footers[0].Content)
	testutils.Equal(t, "second thing", logs1[0].Footers[1].Content)
	testutils.Equal(t, "third thing", logs1[0].Footers[2].Content)
}

func TestParseBreakingChange2(t *testing.T) {
	// And commit must maintain GH link to PR like #19
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := NewSemanticRepo(gitInMemory)
	gitInMemory.TestCommit(`feat: smth is breaking

Some descreiption to why it is happening
which i am writing in multiline but wihtout two dots thingy
bla bla

BREAKING CHANGE: first thing
	multiline first continued
BREAKING CHANGE: second thing
BREAKING CHANGE: third thing

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# Date:      Sun Oct 29 02:38:37 2023 +0100
#
# On branch master
# Changes to be committed:
#       modified:   README.md
#
# Changes not staged for commit:
#       deleted:    autogit.yml
#`)

	logs1 := gitSemantic.GetChangelogByTag("", true)
	assert.Len(t, logs1, 1)

	assert.Len(t, logs1[0].Footers, 3)

	testutils.Equal(t, "first thing\n\tmultiline first continued", logs1[0].Footers[0].Content)
	testutils.Equal(t, "second thing", logs1[0].Footers[1].Content)
	testutils.Equal(t, "third thing", logs1[0].Footers[2].Content)
}

func TestParseBreakingChange3(t *testing.T) {
	// And commit must maintain GH link to PR like #19
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := NewSemanticRepo(gitInMemory)
	gitInMemory.TestCommit(`feat!: new super feature

BREAKING CHANGE: api for endpoint /status changed to /users/status`)

	logs1 := gitSemantic.GetChangelogByTag("", true)
	assert.Len(t, logs1, 1)

	assert.Len(t, logs1[0].Footers, 1)

	testutils.Equal(t, "api for endpoint /status changed to /users/status", logs1[0].Footers[0].Content)
}

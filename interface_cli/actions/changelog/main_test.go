package changelog

import (
	"autogit/interface_cli/actions/changelog/changelog_types"
	"autogit/settings"
	"autogit/settings/testutils"
	_ "autogit/settings/testutils/autouse"
	"autogit/settings/types"

	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"autogit/semanticgit/semver/semvertype"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FixtureGitSemantic(t *testing.T) (*git.Repository, *semanticgit.SemanticGit, settings.ConfigScheme) {
	gitInMemory := git.NewRepoTestInMemory()
	gitSemantic := semanticgit.NewSemanticRepo(gitInMemory)

	gitInMemory.TestCommit("feat: test")

	testutils.Equal(t, "v0.0.0", gitSemantic.GetCurrentVersion().ToString())

	gitInMemory.TestCreateTag("v0.0.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: test2")

	testutils.Equal(t, "v0.0.1", gitSemantic.GetCurrentVersion().ToString())
	return gitInMemory, gitSemantic, settings.GetConfig()
}

func TestPrepare(t *testing.T) {
	_, gitSemantic, config := FixtureGitSemantic(t)

	// Just for debug
	rendered := NewChangelog(gitSemantic, semvertype.OptionsSemVer{}, config, types.TagName(""))
	_ = rendered
	rendered = NewChangelog(gitSemantic, semvertype.OptionsSemVer{}, config, types.TagName("v0.0.1"))
	_ = rendered
}

func TestRender(t *testing.T) {
	_, gitSemantic, config := FixtureGitSemantic(t)

	rendered := NewChangelog(gitSemantic, semvertype.OptionsSemVer{}, config, types.TagName("")).Render()
	assert.Contains(t, rendered, "v0.1.0")

	// historing render
	rendered = NewChangelog(gitSemantic, semvertype.OptionsSemVer{}, config, types.TagName("v0.0.1")).Render()
	assert.Contains(t, rendered, "v0.0.1")
}

func CountCommitsInChangelog(section changelog_types.ChangelogSection, changelog changelogVars) int {
	changeloged_merge_commits_count := 0
	section_group, ok := changelog.SemverGroups[section]
	if !ok {
		return changeloged_merge_commits_count
	}

	for _, commit_group := range section_group.CommitTypeGroups {
		changeloged_merge_commits_count += len(commit_group.NoScopeCommits)
		for _, commits := range commit_group.ScopedCommits {
			changeloged_merge_commits_count += len(commits)
		}
	}
	return changeloged_merge_commits_count
}

func TestMergeCommitsInChangelog(t *testing.T) {
	gitInMemory, gitSemantic, config := FixtureGitSemantic(t)

	gitInMemory.TestCommit("merge: pull request #5 from some branch")
	gitInMemory.TestCommit("Merge pull request #5 from some branch") // Github standard msg
	gitInMemory.TestCommit("merge: unknown pull request from some branch")

	config.Changelog.MergeCommits.MustHaveLinkedPR = true
	config.Changelog.MergeCommits.RedirectMergingCommits = false
	changelog := NewChangelog(gitSemantic, semvertype.OptionsSemVer{}, config, types.TagName(""))
	rendered := changelog.Render()
	assert.Contains(t, rendered, "v0.1.0", "for MustHaveLinkedPR=true, RedirectMergingCommits=false")
	assert.Equal(t, 2, CountCommitsInChangelog(MergeCommits, changelog))

	config.Changelog.MergeCommits.MustHaveLinkedPR = false
	config.Changelog.MergeCommits.RedirectMergingCommits = false
	changelog = NewChangelog(gitSemantic, semvertype.OptionsSemVer{}, config, types.TagName(""))
	rendered = changelog.Render()
	assert.Contains(t, rendered, "v0.1.0", "for MustHaveLinkedPR=false, RedirectMergingCommits=false")
	assert.Equal(t, 3, CountCommitsInChangelog(MergeCommits, changelog))
	assert.Equal(t, 1, CountCommitsInChangelog(SemVerMinor, changelog))
	assert.Equal(t, 0, CountCommitsInChangelog(SemVerPatch, changelog))

	config.Changelog.MergeCommits.MustHaveLinkedPR = true
	config.Changelog.MergeCommits.RedirectMergingCommits = true
	gitInMemory.TestCommit("merge: pull request #1 from feat/branch")
	gitInMemory.TestCommit("merge: pull request #2 from fix/branch")
	gitInMemory.TestCommit("merge: pull request #3 from fix/branch")
	changelog = NewChangelog(gitSemantic, semvertype.OptionsSemVer{}, config, types.TagName(""))
	rendered = changelog.Render()
	assert.Contains(t, rendered, "v0.1.0", "for MustHaveLinkedPR=false, RedirectMergingCommits=false")

	assert.Equal(t, 2, CountCommitsInChangelog(MergeCommits, changelog), "redirect. Expected 3 merge")
	assert.Equal(t, 2, CountCommitsInChangelog(SemVerMinor, changelog), "redirect. Expected 2 minor")
	assert.Equal(t, 2, CountCommitsInChangelog(SemVerPatch, changelog), "redirect. Expected 2 patch")

	config.Changelog.MergeCommits.MustHaveLinkedPR = false
	config.Changelog.MergeCommits.RedirectMergingCommits = true
	changelog = NewChangelog(gitSemantic, semvertype.OptionsSemVer{}, config, types.TagName(""))
	rendered = changelog.Render()
	assert.Contains(t, rendered, "v0.1.0", "for MustHaveLinkedPR=false, RedirectMergingCommits=false")

	assert.Equal(t, 3, CountCommitsInChangelog(MergeCommits, changelog), "redirect. Expected 3 merge")
	assert.Equal(t, 2, CountCommitsInChangelog(SemVerMinor, changelog), "redirect. Expected 2 minor")
	assert.Equal(t, 2, CountCommitsInChangelog(SemVerPatch, changelog), "redirect. Expected 2 patch")
}

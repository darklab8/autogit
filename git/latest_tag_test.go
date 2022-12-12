package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestTag(t *testing.T) {
	repo := (&Repository{}).NewRepoTest()
	repo.TestCommit("feat: test")
	repo.TestCommit("feat: test5")
	assert.Equal(t, "", repo.GetLatestTagString())

	repo.TestCreateTag("v0.0.1", repo.TestCommit("fix: thing"))
	assert.Equal(t, "v0.0.1", repo.GetLatestTagString())

	repo.TestCommit("feat(api): test")
	assert.Equal(t, "v0.0.1", repo.GetLatestTagString())

	repo.TestCreateTag("v0.0.2", repo.TestCommit("feat(api): test2"))
	assert.Equal(t, "v0.0.2", repo.GetLatestTagString())
	repo.TestCommit("fix: test1")
	repo.TestCommit("fix: test2")
	repo.TestCommit("fix: test3")
	assert.Equal(t, "v0.0.2", repo.GetLatestTagString())
}

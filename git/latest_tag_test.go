package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestTag(t *testing.T) {
	repo := (&TestRepository{}).New()
	repo.Commit("feat: test")
	repo.Commit("feat: test5")
	assert.Equal(t, "", repo.GetLatestTagString())

	repo.CreateTag("v0.0.1", repo.Commit("fix: thing"))
	assert.Equal(t, "v0.0.1", repo.GetLatestTagString())

	repo.Commit("feat(api): test")
	assert.Equal(t, "v0.0.1", repo.GetLatestTagString())

	repo.CreateTag("v0.0.2", repo.Commit("feat(api): test2"))
	assert.Equal(t, "v0.0.2", repo.GetLatestTagString())
	repo.Commit("fix: test1")
	repo.Commit("fix: test2")
	repo.Commit("fix: test3")
	assert.Equal(t, "v0.0.2", repo.GetLatestTagString())
}

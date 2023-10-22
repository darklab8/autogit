package git

import (
	"autogit/settings/testutils"
	"autogit/settings/types"
	"testing"
)

func (r *Repository) GetLatestTagString() types.TagName {
	var return_tag Tag

	r.ForeachTag(func(tag Tag) bool {
		return_tag = tag
		return true
	})

	return return_tag.Name
}

func TestGetLatestTag(t *testing.T) {
	repo := NewRepoTestInMemory()
	repo.TestCommit("feat: test")
	repo.TestCommit("feat: test5")
	testutils.Equal(t, types.TagName(""), repo.GetLatestTagString())

	repo.TestCreateTag("v0.0.1", repo.TestCommit("fix: thing"))
	testutils.Equal(t, types.TagName("v0.0.1"), repo.GetLatestTagString())

	repo.TestCommit("feat(api): test")
	testutils.Equal(t, types.TagName("v0.0.1"), repo.GetLatestTagString())

	repo.TestCreateTag("v0.0.2", repo.TestCommit("feat(api): test2"))
	testutils.Equal(t, types.TagName("v0.0.2"), repo.GetLatestTagString())
	repo.TestCommit("fix: test1")
	repo.TestCommit("fix: test2")
	repo.TestCommit("fix: test3")
	testutils.Equal(t, types.TagName("v0.0.2"), repo.GetLatestTagString())
	repo.TestCreateTag("v0.0.3", repo.TestCommit("feat(api): test2"))
	testutils.Equal(t, types.TagName("v0.0.3"), repo.GetLatestTagString())
}

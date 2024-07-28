package git

import (
	"fmt"
	"testing"

	"github.com/darklab8/autogit/v2/settings/testutils"
)

func TestGitRepo(t *testing.T) {
	repo := NewRepoTestInMemory()
	repo.TestCommit("feat: test")
	repo.TestCommit("feat: test3")
	repo.TestCommit("feat: test5")
	repo.TestCreateTag("v0.0.1", repo.TestCommit("fix: thing"))
	repo.TestCommit("feat(api): test")
	repo.TestCreateTag("v0.0.2", repo.TestCommit("feat(api): test2"))
	repo.TestCommit("fix: test1")
	repo.TestCommit("fix: test2")
	repo.TestCommit("fix: test3")

	tags := repo.getUnorderedTags()
	fmt.Printf("tags=%v\n", tags)
	testutils.Equal(t, 2, len(tags))
}

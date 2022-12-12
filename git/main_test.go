package git

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitRepo(t *testing.T) {
	repo := (&Repository{}).NewRepoTest()
	repo.TestCommit("feat: test")
	repo.TestCommit("feat: test3")
	repo.TestCommit("feat: test5")
	repo.TestCreateTag("v0.0.1", repo.TestCommit("fix: thing"))
	repo.TestCommit("feat(api): test")
	repo.TestCreateTag("v0.0.2", repo.TestCommit("feat(api): test2"))
	repo.TestCommit("fix: test1")
	repo.TestCommit("fix: test2")
	repo.TestCommit("fix: test3")

	tags := repo.GetTags()
	fmt.Printf("tags=%v\n", tags)
	assert.Equal(t, 2, len(tags))

	logs1 := repo.TestGetChangelogByTag("")
	assert.Len(t, logs1, 3)

	logs2 := repo.TestGetChangelogByTag("v0.0.2")
	assert.Len(t, logs2, 2)

	logs3 := repo.TestGetChangelogByTag("v0.0.1")
	assert.Len(t, logs3, 4)
}

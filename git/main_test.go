package git

import (
	"autogit/utils"
	"fmt"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/stretchr/testify/assert"
)

type TestRepository struct {
	Repository
	wt     *git.Worktree
	author *object.Signature
}

func (r *TestRepository) New() *TestRepository {
	var err error
	fs := memfs.New()

	r.repo, err = git.Init(memory.NewStorage(), fs)
	CheckIfError(err)

	r.wt, err = r.repo.Worktree()

	r.author = &object.Signature{Name: "abc", Email: "abc@example.com"}
	return r
}

func (r *TestRepository) Commit(msg string) plumbing.Hash {
	file := (&utils.File{Filepath: "testfile.txt"}).CreateToWriteF()
	defer file.Close()
	file.WritelnF(utils.TokenHex(10))

	r.wt.Add("testfile.txt")
	hash, err := r.wt.Commit(msg, &git.CommitOptions{Author: r.author})
	CheckIfError(err)
	return hash
}

func (r *TestRepository) CreateTag(name string, hash plumbing.Hash) {
	ref, err := r.repo.CreateTag(name, hash, &git.CreateTagOptions{Tagger: r.author, Message: "123"})
	fmt.Printf("CreateTag=%v,%v\n", ref, err)
}

func (r *Repository) GetChangelogByTag(tagName string) []Log {
	if tagName == "" {
		return r.GetLogs(HEAD_Hash)
	}

	tag_ref, _ := r.repo.Tag(tagName)
	tag_obj, _ := r.repo.TagObject(tag_ref.Hash())
	logs := r.GetLogs(tag_obj.Target)
	return logs
}

func TestGitRepo(t *testing.T) {
	repo := (&TestRepository{}).New()
	repo.Commit("feat: test")
	repo.Commit("feat: test3")
	repo.Commit("feat: test5")
	repo.CreateTag("v0.0.1", repo.Commit("fix: thing"))
	repo.Commit("feat(api): test")
	repo.CreateTag("v0.0.2", repo.Commit("feat(api): test2"))
	repo.Commit("fix: test1")
	repo.Commit("fix: test2")
	repo.Commit("fix: test3")

	tags := repo.GetTags()
	fmt.Printf("tags=%v\n", tags)
	assert.Equal(t, 2, len(tags))

	logs1 := repo.GetChangelogByTag("")
	assert.Len(t, logs1, 3)

	logs2 := repo.GetChangelogByTag("v0.0.2")
	assert.Len(t, logs2, 2)

	logs3 := repo.GetChangelogByTag("v0.0.1")
	assert.Len(t, logs3, 4)
}

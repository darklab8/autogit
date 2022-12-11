package git

import (
	"autogit/utils"
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

func TestSaveRecycleParams(t *testing.T) {
	repo := (&TestRepository{}).New()
	repo.Commit("feat: test")
	repo.Commit("fix: thing")
	repo.Commit("feat{api): test")
	repo.GetLogs()

	assert.True(t, true)
}

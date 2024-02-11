package git

import (
	"fmt"

	"github.com/darklab8/autogit/settings/logus"

	"github.com/darklab8/go-utils/goutils/utils"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func NewRepoTestInMemory() *Repository {
	r := &Repository{}
	var err error
	fs := memfs.New()

	r.repo, err = git.Init(memory.NewStorage(), fs)
	logus.Log.CheckFatal(err, "failed git init")

	r.wt, err = r.repo.Worktree()
	logus.Log.CheckFatal(err, "failed to get worktree of repo")

	r.author = &object.Signature{Name: "abc", Email: "abc@example.com"}
	return r
}

func (r *Repository) TestCommit(msg string) plumbing.Hash {
	utils.NewWriteFile("testfile.txt", func(file *utils.FileWrite) {
		file.WritelnF(utils.TokenHex(10))
	})

	r.wt.Add("testfile.txt")
	hash, err := r.wt.Commit(msg, &git.CommitOptions{Author: r.author})
	logus.Log.CheckFatal(err, "unable to form commit")
	return hash
}

func (r *Repository) TestCreateTag(name string, hash plumbing.Hash) {
	ref, err := r.repo.CreateTag(name, hash, &git.CreateTagOptions{Tagger: r.author, Message: "123"})
	fmt.Printf("CreateTag=%v,%v\n", ref, err)
}

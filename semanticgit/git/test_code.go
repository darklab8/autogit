package git

import (
	"autogit/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func (r *Repository) TestNewRepo() *Repository {
	var err error
	fs := memfs.New()

	r.repo, err = git.Init(memory.NewStorage(), fs)
	CheckIfError(err)

	r.wt, err = r.repo.Worktree()

	r.author = &object.Signature{Name: "abc", Email: "abc@example.com"}
	return r
}

func (r *Repository) TestNewRepoIntegration() *Repository {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err, "unable to get workdir")
	}
	r.repo, err = git.PlainOpen(filepath.Dir(filepath.Dir(path)))
	if err != nil {
		log.Fatal(err, "unable to open git")
	}
	return r
}

func (r *Repository) TestCommit(msg string) plumbing.Hash {
	file := (&utils.File{Filepath: "testfile.txt"}).CreateToWriteF()
	defer file.Close()
	file.WritelnF(utils.TokenHex(10))

	r.wt.Add("testfile.txt")
	hash, err := r.wt.Commit(msg, &git.CommitOptions{Author: r.author})
	CheckIfError(err)
	return hash
}

func (r *Repository) TestCreateTag(name string, hash plumbing.Hash) {
	ref, err := r.repo.CreateTag(name, hash, &git.CreateTagOptions{Tagger: r.author, Message: "123"})
	fmt.Printf("CreateTag=%v,%v\n", ref, err)
}

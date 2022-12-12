package git

import (
	"autogit/utils"
	"fmt"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func (r *Repository) NewRepoTest() *Repository {
	var err error
	fs := memfs.New()

	r.repo, err = git.Init(memory.NewStorage(), fs)
	CheckIfError(err)

	r.wt, err = r.repo.Worktree()

	r.author = &object.Signature{Name: "abc", Email: "abc@example.com"}
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

func (r *Repository) TestGetChangelogByTag(tagName string) []Log {
	if tagName == "" {
		return r.GetLogs(HEAD_Hash)
	}

	tag_ref, _ := r.repo.Tag(tagName)
	tag_obj, _ := r.repo.TagObject(tag_ref.Hash())
	logs := r.GetLogs(tag_obj.Target)
	return logs
}

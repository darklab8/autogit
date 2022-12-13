// Like git module. But wrapper to one place
package git

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CheckIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Repository struct {
	repo   *git.Repository
	wt     *git.Worktree
	author *object.Signature
}

func (r *Repository) NewRepoIntegration() *Repository {
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

func (r *Repository) NewRepoInWorkDir() *Repository {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err, "unable to get workdir")
	}
	r.repo, err = git.PlainOpen(path)
	if err != nil {
		log.Fatal(err, "unable to open git")
	}
	return r
}

type Log struct {
	Hash plumbing.Hash
	Msg  string
}

var HEAD_Hash plumbing.Hash

func (r *Repository) GetLatestTag(skipLatestCommit bool) Tag {
	ref, err := r.repo.Head()
	CheckIfError(err)
	From := ref.Hash()

	// ... retrieves the commit history
	cIter, err := r.repo.Log(&git.LogOptions{From: From})
	CheckIfError(err)

	tags := r.GetTags()
	// ... just iterates over the commits, printing it
	c, _ := cIter.Next()
	firstCommit := c
	for ; c != nil; c, _ = cIter.Next() {
		// iterating until next tag
		for _, tag := range tags {
			if skipLatestCommit && tag.Hash == firstCommit.Hash {
				continue
			} else if tag.Hash == c.Hash {
				return tag
			}
		}

	}
	CheckIfError(err)

	return Tag{}
}

func (r *Repository) GetLatestTagString(skipLatestCommit bool) string {
	return r.GetLatestTag(skipLatestCommit).Name
}

func (r *Repository) GetLogs(From plumbing.Hash) []Log {
	var logs []Log
	// retrieves the branch pointed by HEAD
	if From.IsZero() {
		var err error
		ref, err := r.repo.Head()
		CheckIfError(err)
		From = ref.Hash()
	}

	// get the commit object, pointed by ref
	// commit, err := r.CommitObject(ref.Hash())

	// ... retrieves the commit history
	cIter, err := r.repo.Log(&git.LogOptions{From: From})
	CheckIfError(err)

	tags := r.GetTags()
	// ... just iterates over the commits, printing it
	c, _ := cIter.Next()
	for ; c != nil; c, _ = cIter.Next() {

		// iterating until next tag
		if From != c.Hash {
			for _, tag := range tags {
				if tag.Hash == c.Hash {
					return logs
				}
			}

		}

		logs = append(logs, Log{Hash: c.Hash, Msg: c.Message})
	}
	CheckIfError(err)

	return logs
}

type Tag struct {
	Hash plumbing.Hash
	Ref  *plumbing.Reference
	Name string
}

// brings tags, from latest to new ones
func (r *Repository) GetTags() []Tag {
	var results []Tag
	iter, err := r.repo.Tags()
	CheckIfError(err)

	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		if !ref.Name().IsTag() {
			return nil
		}
		tag, err := r.repo.Tag(ref.Name().Short())
		if err != nil {
			log.Fatal("failed to get tag ", ref.Name())
		}

		tag_obj, err := r.repo.TagObject(ref.Hash())
		if err == nil {
			results = append(results, Tag{Hash: tag_obj.Target, Name: tag_obj.Name, Ref: ref})
			return nil
		}

		results = append(results, Tag{Hash: tag.Hash(), Name: tag.Name().Short(), Ref: ref})
		return nil
	}); err != nil {
		CheckIfError(err)
	}

	return results
}

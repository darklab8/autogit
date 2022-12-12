// Like git module. But wrapper to one place
package git

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CheckIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Repository struct {
	repo *git.Repository
}

func (r *Repository) GetRepoInWorkDir() *Repository {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Getwd=%s", path) // for example /home/user

	r.repo, err = git.PlainOpen(path)
	return r
}

type Log struct {
	Hash plumbing.Hash
	Msg  string
}

var HEAD_Hash plumbing.Hash

func (r *Repository) GetLatestTag() Tag {
	ref, err := r.repo.Head()
	CheckIfError(err)
	From := ref.Hash()

	// ... retrieves the commit history
	cIter, err := r.repo.Log(&git.LogOptions{From: From})
	CheckIfError(err)

	tags := r.GetTags()
	// ... just iterates over the commits, printing it
	c, _ := cIter.Next()
	for ; c != nil; c, _ = cIter.Next() {

		// iterating until next tag
		for _, tag := range tags {
			if tag.Target == c.Hash {
				return Tag{Hash: tag.Hash, Name: tag.Name, Message: tag.Message, Ref: ref, Target: tag.Target}
			}
		}

	}
	CheckIfError(err)

	return Tag{Name: "v0.0.0", Message: "default"}
}

func (r *Repository) GetLatestTagString() string {
	return r.GetLatestTag().Name
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
				if tag.Target == c.Hash {
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
	Hash    plumbing.Hash
	Target  plumbing.Hash
	Ref     *plumbing.Reference
	Name    string
	Message string
}

// brings tags, from latest to new ones
func (r *Repository) GetTags() []Tag {
	var results []Tag
	iter, err := r.repo.Tags()
	CheckIfError(err)

	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		obj, err := r.repo.TagObject(ref.Hash())
		// fmt.Printf("%s-%s-%s\n", obj.Name, obj.Hash, obj.Message)
		switch err {
		case nil:
			// Tag object present
			results = append(results, Tag{Hash: obj.Hash, Name: obj.Name, Message: obj.Message, Ref: ref, Target: obj.Target})
			ref.Name()
			ref.Target()
			r.repo.Object(plumbing.TagObject, ref.Hash())
		case plumbing.ErrObjectNotFound:
			// Not a tag object
		default:
			// Some other error
			return err
		}
		return nil
	}); err != nil {
		CheckIfError(err)
	}

	return results
}

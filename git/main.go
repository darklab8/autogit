// Like git module. But wrapper to one place
package git

import (
	"fmt"
	"log"
	"os"

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

func (r *Repository) GetLogs() {

	// retrieves the branch pointed by HEAD
	ref, err := r.repo.Head()

	// get the commit object, pointed by ref
	// commit, err := r.CommitObject(ref.Hash())

	// ... retrieves the commit history
	cIter, err := r.repo.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})
	CheckIfError(err)
}

type Tag struct {
	Hash    plumbing.Hash
	Name    string
	Message string
}

func (r *Repository) GetTags() []Tag {
	var results []Tag
	iter, err := r.repo.Tags()
	CheckIfError(err)

	// ref, err := iter.Next()
	// obj, err := r.repo.TagObject(ref.Hash())
	// results = append(results, Tag{Hash: obj.Hash, Name: obj.Name, Message: obj.Message})

	// if err := iter.ForEach(func(ref *plumbing.Reference) error {
	// 	obj, err := r.repo.TagObject(ref.Hash())
	// 	results = append(results, Tag{Hash: obj.Hash, Name: obj.Name, Message: obj.Message})
	// 	return err
	// }); err != nil {
	// 	CheckIfError(err)
	// }

	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		obj, err := r.repo.TagObject(ref.Hash())
		fmt.Printf("%s-%s-%s\n", obj.Name, obj.Hash, obj.Message)
		switch err {
		case nil:
			// Tag object present
			results = append(results, Tag{Hash: obj.Hash, Name: obj.Name, Message: obj.Message})
			fmt.Printf("tag=%v\n", obj)
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

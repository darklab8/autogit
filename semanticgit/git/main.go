// Like git module. But wrapper to one place
package git

import (
	"autogit/settings"
	"autogit/utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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

func (r *Repository) GetLatestCommitHash() plumbing.Hash {
	ref, err := r.repo.Head()
	CheckIfError(err)
	return ref.Hash()
}

func (r *Repository) ForeachTag(callback func(tag Tag) bool) {
	From := r.GetLatestCommitHash()

	// ... retrieves the commit history
	cIter, err := r.repo.Log(&git.LogOptions{From: From})
	CheckIfError(err)

	tags := r.getUnorderedTags()
	// ... just iterates over the commits, printing it
	c, _ := cIter.Next()
	for ; c != nil; c, _ = cIter.Next() {
		// iterating until next tag
		for _, tag := range tags {
			if tag.Hash == c.Hash {
				shouldWeStop := callback(tag)
				if shouldWeStop {
					return
				}
			}
		}

	}
	CheckIfError(err)
}

func (r *Repository) ForeachLog(From plumbing.Hash, callback func(log Log) bool) {
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

	// ... just iterates over the commits, printing it
	c, _ := cIter.Next()
	for ; c != nil; c, _ = cIter.Next() {
		shouldWeStop := callback(Log{Hash: c.Hash, Msg: c.Message})
		if shouldWeStop {
			return
		}
	}
	CheckIfError(err)
}

type Tag struct {
	Hash plumbing.Hash
	Ref  *plumbing.Reference
	Name string
}

// brings tags, from latest to new ones
func (r *Repository) getUnorderedTags() []Tag {
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

func (r *Repository) getHashByTagName(tagName string) plumbing.Hash {
	if tagName == "" {
		return HEAD_Hash
	}

	tag_ref, _ := r.repo.Tag(tagName)
	tag_obj, err := r.repo.TagObject(tag_ref.Hash())
	if err == nil {
		return tag_obj.Target
	}
	return tag_ref.Hash()
}

func (r *Repository) GetLogsFromTag(tagName string, callback func(log Log) bool) {
	FromHash := r.getHashByTagName(tagName)

	r.ForeachLog(FromHash, func(log Log) bool {
		return callback(log)
	})
}

func (r *Repository) CreateTag(name string, msg string) {
	hash, err := r.repo.Head()
	utils.CheckFatal(err)
	ref, err := r.repo.CreateTag(name, hash.Hash(), &git.CreateTagOptions{Message: msg})
	fmt.Printf("CreateTag=%v,%v\n", ref, err)
}

const defaultRemoteName = "origin"

func (r *Repository) PushTag(name string) {
	var publicKey *ssh.PublicKeys
	sshPath := filepath.Join(os.Getenv("HOME"), ".ssh", settings.Config.Git.SSHPath)
	sshKey, _ := ioutil.ReadFile(sshPath)
	publicKey, keyError := ssh.NewPublicKeys("git", []byte(sshKey), "")
	utils.CheckFatal(keyError)

	refs := []config.RefSpec{
		config.RefSpec("+refs/tags/" + name + ":refs/tags/" + name),
	}
	utils.CheckFatal(refs[0].Validate(), "failed to validate push tag")
	err := r.repo.Push(&git.PushOptions{RemoteName: defaultRemoteName, Auth: publicKey, RefSpecs: refs, Progress: os.Stdout})
	utils.CheckFatal(err, "failed to push")
	fmt.Printf("PushTag=%v\n", err)
}

func (r *Repository) HookEnabled(enabled bool) {
	hooksPathkey := "hooksPath"
	cfg, err := r.repo.Config()
	utils.CheckFatal(err, "failed to read config")

	if enabled {
		cfg.Raw.Section("core").SetOption(hooksPathkey, ".git-hooks")
	} else {
		cfg.Raw.Section("core").RemoveOption(hooksPathkey)
	}
	r.repo.SetConfig(cfg)
	utils.CheckFatal(err, "failed to write config")
}

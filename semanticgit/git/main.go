// Like git module. But wrapper to one place
package git

import (
	"autogit/settings/logus"
	"autogit/settings/types"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type SshPath string

type Repository struct {
	repo    *git.Repository
	wt      *git.Worktree
	author  *object.Signature
	sshPath SshPath
}

func (r *Repository) NewRepoInWorkDir(sshPath SshPath) *Repository {
	r.sshPath = sshPath

	path, err := os.Getwd()
	logus.CheckFatal(err, "unable to get workdir")
	r.repo, err = git.PlainOpen(path)
	logus.CheckFatal(err, "unable to open git")
	return r
}

type Log struct {
	Hash plumbing.Hash
	Msg  types.CommitMessage
}

var HEAD_Hash plumbing.Hash

func (r *Repository) GetLatestCommitHash() plumbing.Hash {
	ref, err := r.repo.Head()
	logus.CheckFatal(err, "unabled to get latest Head commit")
	return ref.Hash()
}

func (r *Repository) ForeachTag(callback func(tag Tag) bool) {
	From := r.GetLatestCommitHash()

	// ... retrieves the commit history
	cIter, err := r.repo.Log(&git.LogOptions{From: From})
	logus.CheckFatal(err, "unable to get git log")

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
}

func (r *Repository) ForeachLog(From plumbing.Hash, callback func(log Log) bool) {
	// retrieves the branch pointed by HEAD
	if From.IsZero() {
		var err error
		ref, err := r.repo.Head()
		logus.CheckFatal(err, "unable getting Head commit")
		From = ref.Hash()
	}

	// get the commit object, pointed by ref
	// commit, err := r.CommitObject(ref.Hash())

	// ... retrieves the commit history
	cIter, err := r.repo.Log(&git.LogOptions{From: From})
	logus.CheckFatal(err, "unable getting git log")

	// ... just iterates over the commits, printing it
	c, _ := cIter.Next()
	for ; c != nil; c, _ = cIter.Next() {
		shouldWeStop := callback(Log{Hash: c.Hash, Msg: types.CommitMessage(c.Message)})
		if shouldWeStop {
			return
		}
	}
}

type Tag struct {
	Hash plumbing.Hash
	Ref  *plumbing.Reference
	Name types.TagName
}

// brings tags, from latest to new ones
func (r *Repository) getUnorderedTags() []Tag {
	var results []Tag
	iter, err := r.repo.Tags()
	logus.CheckFatal(err, "unable to get repository tags")

	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		if !ref.Name().IsTag() {
			return nil
		}
		tag, err := r.repo.Tag(ref.Name().Short())
		tag_name := types.TagName(tag.Name())
		if err != nil {
			logus.Fatal("failed to get tag ", logus.TagName(tag_name))
		}

		tag_obj, err := r.repo.TagObject(ref.Hash())
		if err == nil {
			results = append(results, Tag{Hash: tag_obj.Target, Name: types.TagName(tag_obj.Name), Ref: ref})
			return nil
		}

		results = append(results, Tag{Hash: tag.Hash(), Name: types.TagName(tag.Name().Short()), Ref: ref})
		return nil
	}); err != nil {
		logus.CheckFatal(err, "failed iterating repository refs")
	}

	return results
}

func (r *Repository) getHashByTagName(tagName types.TagName) plumbing.Hash {
	if tagName == "" {
		return HEAD_Hash
	}

	tag_ref, _ := r.repo.Tag(string(tagName))
	tag_obj, err := r.repo.TagObject(tag_ref.Hash())
	if err == nil {
		return tag_obj.Target
	}
	return tag_ref.Hash()
}

func (r *Repository) GetLogsFromTag(tagName types.TagName, callback func(log Log) bool) {
	FromHash := r.getHashByTagName(tagName)

	r.ForeachLog(FromHash, func(log Log) bool {
		return callback(log)
	})
}

func (r *Repository) CreateTag(name types.TagName, msg string) {
	hash, err := r.repo.Head()
	logus.CheckFatal(err, "failed getting Head commit")
	ref, err := r.repo.CreateTag(string(name), hash.Hash(), &git.CreateTagOptions{Message: msg})
	fmt.Printf("CreateTag=%v,%v\n", ref, err)
}

const defaultRemoteName = "origin"

func (r *Repository) PushTag(name types.TagName) {
	var publicKey *ssh.PublicKeys
	sshPath := filepath.Join(os.Getenv("HOME"), ".ssh", string(r.sshPath))
	sshKey, _ := os.ReadFile(sshPath)
	publicKey, keyError := ssh.NewPublicKeys("git", []byte(sshKey), "")
	logus.CheckFatal(keyError, "failed initializing git ssh keys")

	refs := []config.RefSpec{
		config.RefSpec("+refs/tags/" + name + ":refs/tags/" + name),
	}
	logus.CheckFatal(refs[0].Validate(), "failed to validate push tag")
	err := r.repo.Push(&git.PushOptions{RemoteName: defaultRemoteName, Auth: publicKey, RefSpecs: refs, Progress: os.Stdout})
	logus.CheckFatal(err, "failed to push")
	fmt.Printf("PushTag=%v\n", err)
}

func (r *Repository) HookEnabled(enabled bool) {
	hooksPathkey := "hooksPath"
	cfg, err := r.repo.Config()
	logus.CheckFatal(err, "failed to read config")

	if enabled {
		cfg.Raw.Section("core").SetOption(hooksPathkey, ".git-hooks")
	} else {
		cfg.Raw.Section("core").RemoveOption(hooksPathkey)
	}
	r.repo.SetConfig(cfg)
	logus.CheckFatal(err, "failed to write config")
}

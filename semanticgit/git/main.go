// Like git module. But wrapper to one place
package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/darklab8/autogit/semanticgit/git/gitraw"
	"github.com/darklab8/autogit/settings"
	"github.com/darklab8/autogit/settings/envs"
	"github.com/darklab8/autogit/settings/logus"
	"github.com/darklab8/autogit/settings/types"
	"github.com/darklab8/go-typelog/typelog"

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

func NewRepoInWorkDir(sshPath SshPath) *Repository {
	r := &Repository{}
	r.sshPath = sshPath
	r.repo = gitraw.NewGitRepo()
	return r
}

type Log struct {
	Hash plumbing.Hash
	Msg  types.CommitOriginalMsg
}

var HEAD_Hash plumbing.Hash

func (r *Repository) GetLatestCommitHash() plumbing.Hash {
	ref, err := r.repo.Head()
	logus.Log.CheckFatal(err, "unabled to get latest Head commit")
	return ref.Hash()
}

type ShouldWeStopIteration bool

func (r *Repository) ForeachTag(callback func(tag Tag) ShouldWeStopIteration) {
	From := r.GetLatestCommitHash()

	// ... retrieves the commit history
	cIter, err := r.repo.Log(&git.LogOptions{From: From})
	logus.Log.CheckFatal(err, "unable to get git log")

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

func (r *Repository) ForeachLog(From plumbing.Hash, callback func(log Log) ShouldWeStopIteration) {
	logus.Log.Debug(fmt.Sprintf("ForeachLog. From1=%v", From))
	// retrieves the branch pointed by HEAD
	if From.IsZero() {
		var err error
		ref, err := r.repo.Head()
		logus.Log.CheckFatal(err, "unable getting Head commit")
		From = ref.Hash()
	}
	logus.Log.Debug(fmt.Sprintf("ForeachLog. From2=%v", From))

	// get the commit object, pointed by ref
	// commit, err := r.CommitObject(ref.Hash())

	// ... retrieves the commit history
	cIter, err := r.repo.Log(
		&git.LogOptions{
			From:  From,
			Order: git.LogOrderCommitterTime, // necessary to see commits between merging commits
		},
	)
	logus.Log.CheckFatal(err, "unable getting git log")

	// ... just iterates over the commits, printing it
	c, _ := cIter.Next()
	for ; c != nil; c, _ = cIter.Next() {
		msg := types.CommitOriginalMsg(c.Message)
		// autogit_logus.Log.Debug("ForeachLog retrieved", logus_core.CommitMessage(msg), logus_core.CommitHash(c.Hash))
		shouldWeStop := callback(Log{Hash: c.Hash, Msg: msg})
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
	logus.Log.CheckFatal(err, "unable to get repository tags")

	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		if !ref.Name().IsTag() {
			return nil
		}
		tag, err := r.repo.Tag(ref.Name().Short())
		tag_name := types.TagName(tag.Name())
		if err != nil {
			logus.Log.Fatal("failed to get tag ", logus.TagName(tag_name))
		}

		tag_obj, err := r.repo.TagObject(ref.Hash())
		if err == nil {
			results = append(results, Tag{Hash: tag_obj.Target, Name: types.TagName(tag_obj.Name), Ref: ref})
			return nil
		}

		results = append(results, Tag{Hash: tag.Hash(), Name: types.TagName(tag.Name().Short()), Ref: ref})
		return nil
	}); err != nil {
		logus.Log.CheckFatal(err, "failed iterating repository refs")
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

func (r *Repository) GetLogsFromTag(tagName types.TagName, callback func(log Log) ShouldWeStopIteration) {
	FromHash := r.getHashByTagName(tagName)
	logus.Log.Debug("GetLogsFromTag is called", logus.TagName(tagName))

	r.ForeachLog(FromHash, func(log Log) ShouldWeStopIteration {
		return callback(log)
	})
}

func (r *Repository) CreateTag(name types.TagName, msg string) {
	hash, err := r.repo.Head()
	logus.Log.CheckFatal(err, "failed getting Head commit")
	ref, err := r.repo.CreateTag(string(name), hash.Hash(), &git.CreateTagOptions{Message: msg})
	fmt.Printf("CreateTag=%v,%v\n", ref, err)
}

const defaultRemoteName = "origin"

func (r *Repository) PushTag(name types.TagName) {
	var publicKey *ssh.PublicKeys
	sshPath := filepath.Join(string(envs.PathUserHome), ".ssh", string(r.sshPath))
	sshKey, _ := os.ReadFile(sshPath)
	publicKey, keyError := ssh.NewPublicKeys("git", []byte(sshKey), "")
	logus.Log.CheckFatal(keyError, "failed initializing git ssh keys", typelog.String("key_path", string(sshPath)))

	refs := []config.RefSpec{
		config.RefSpec("+refs/tags/" + name + ":refs/tags/" + name),
	}
	logus.Log.CheckFatal(refs[0].Validate(), "failed to validate push tag")
	err := r.repo.Push(&git.PushOptions{RemoteName: defaultRemoteName, Auth: publicKey, RefSpecs: refs, Progress: os.Stdout})
	logus.Log.CheckFatal(err, "failed to push")
	fmt.Printf("PushTag=%v\n", err)
}

func (r *Repository) HookEnabled(enabled bool) {
	hooksPathkey := "hooksPath"
	cfg, err := r.repo.Config()
	logus.Log.CheckFatal(err, "failed to read config")

	if enabled {
		cfg.Raw.Section("core").SetOption(hooksPathkey, settings.HookFolderName)
	} else {
		cfg.Raw.Section("core").RemoveOption(hooksPathkey)
	}
	r.repo.SetConfig(cfg)
	logus.Log.CheckFatal(err, "failed to write config")
}

/*
Git with applied semantic versioning
and conventional commits to it
*/
package semanticgit

import (
	"autogit/semanticgit/conventionalcommits"
	"autogit/semanticgit/conventionalcommits/conventionalcommitstype"
	"autogit/semanticgit/git"
	"autogit/semanticgit/semver"
	"autogit/semanticgit/semver/semvertype"
	"autogit/settings/logus"
	"autogit/settings/types"
)

type SemanticGit struct {
	git *git.Repository
}

func NewSemanticRepo(gitRepo *git.Repository) *SemanticGit {
	g := &(SemanticGit{})
	g.git = gitRepo
	return g
}

func (g *SemanticGit) GetCurrentVersion() *semvertype.SemVer {
	returned_vers := &semvertype.SemVer{
		Major: 0, Minor: 0, Patch: 0,
		AugmentedSemver: semvertype.AugmentedSemver{Alpha: 0, Beta: 0, Rc: 0},
	}

	latest_hash := g.git.GetLatestCommitHash()
	g.git.ForeachTag(func(tag git.Tag) git.ShouldWeStopIteration {
		vers, err := semver.Parse(tag.Name)

		if err != nil {
			logus.Warn("failed to parse tag=", logus.TagName(tag.Name), logus.OptError(err))
			return git.ShouldWeStopIteration(false)
		}

		if tag.Hash == latest_hash || (vers.Prerelease == "" && tag.Hash == latest_hash) {
			return git.ShouldWeStopIteration(false)
		}

		// Process
		if vers.Rc > returned_vers.Rc {
			returned_vers.Rc = vers.Rc
		}
		if vers.Beta > returned_vers.Beta {
			returned_vers.Beta = vers.Beta
		}
		if vers.Alpha > returned_vers.Alpha {
			returned_vers.Alpha = vers.Alpha
		}

		if vers.Prerelease == "" {
			returned_vers.Major = vers.Major
			returned_vers.Minor = vers.Minor
			returned_vers.Patch = vers.Patch
			return git.ShouldWeStopIteration(true)
		}

		return git.ShouldWeStopIteration(false)
	})

	return returned_vers
}

func (g *SemanticGit) CalculateNextVersion(vers *semvertype.SemVer) *semvertype.SemVer {

	log_records := g.GetChangelogByTag("", false)

	var major_change, minor_change, patch_change bool
	for _, record := range log_records {

		if vers.Options.Publish && vers.Major == 0 {
			major_change = true
		}

		if record.Exclamation {
			if vers.Major != 0 {
				major_change = true
			}
		}

		if record.Type == "feat" {
			minor_change = true
		}

		if record.Type == "fix" {
			patch_change = true
		}
	}

	if major_change {
		vers.Major += 1
		vers.Minor = 0
		vers.Patch = 0
	} else if minor_change {
		vers.Minor += 1
		vers.Patch = 0
	} else if patch_change {
		vers.Patch += 1
	}

	if vers.Options.Alpha {
		vers.Alpha++
	}
	if vers.Options.Beta {
		vers.Beta++
	}
	if vers.Options.Rc {
		vers.Rc++
	}

	// Technically this is the only place where it is set
	// and always from Options to vers
	// preserved in this way to keep package `semver` true to its standard
	// While inputting my all options through Options
	if vers.Options.Build != "" {
		vers.Build = vers.Options.Build
	}
	return vers
}

func (g *SemanticGit) GetNextVersion(semver_options semvertype.OptionsSemVer) *semvertype.SemVer {
	vers := g.GetCurrentVersion()
	vers.Options = semver_options
	vers = g.CalculateNextVersion(vers)

	return vers
}

func (g *SemanticGit) GetChangelogByTag(fromTag types.TagName, enable_warnings bool) []conventionalcommits.ConventionalCommit {
	var result []conventionalcommits.ConventionalCommit

	logus.Debug("semantic git.GetChangelogByTag attempting to get logs", logus.TagName(fromTag))
	g.git.GetLogsFromTag(fromTag, func(log_record git.Log) git.ShouldWeStopIteration {
		parsed_commit, err := conventionalcommits.ParseCommit(log_record.Msg)
		if err != nil {
			if enable_warnings {
				logus.Warn("unable to parse commit with hash=", logus.CommitHash(log_record.Hash), logus.CommitMessage(log_record.Msg))
			}
			logus.Debug("unable to parse commit with hash=", logus.CommitHash(log_record.Hash), logus.CommitMessage(log_record.Msg))
			return git.ShouldWeStopIteration(false)
		}

		// attempt to convert to Semver
		var foundSemver *semvertype.SemVer
		var foundTag git.Tag
		g.git.ForeachTag(func(tag git.Tag) git.ShouldWeStopIteration {
			if log_record.Hash == tag.Hash {
				semver, err := semver.Parse(tag.Name)
				if err != nil {
					return git.ShouldWeStopIteration(false)
				}
				foundSemver = semver
				foundTag = tag
				return git.ShouldWeStopIteration(true)
			}
			return git.ShouldWeStopIteration(false)
		})
		if foundSemver != nil {
			// Get Changelog only until previous stable tag and don't mind first commit
			if foundSemver.Prerelease == "" && fromTag != foundTag.Name && log_record.Hash != g.git.GetLatestCommitHash() {
				logus.Debug("GetChangelogByTag stopping at this commit",
					logus.CommitMessage(log_record.Msg),
					logus.CommitHash(log_record.Hash),
					logus.Semver(foundSemver),
				)
				return git.ShouldWeStopIteration(true)
			}
		}

		if parsed_commit != nil {
			parsed_commit.Hash = conventionalcommitstype.Hash(log_record.Hash.String()[:8])
			result = append(result, *parsed_commit)
		} else {
			logus.Debug("parsed_commit = nil", logus.CommitMessage(log_record.Msg), logus.CommitHash(log_record.Hash))
		}
		return git.ShouldWeStopIteration(false)
	})

	return result
}

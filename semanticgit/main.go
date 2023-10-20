/*
Git with applied semantic versioning
and conventional commits to it
*/
package semanticgit

import (
	"autogit/semanticgit/conventionalcommits"
	"autogit/semanticgit/git"
	"autogit/semanticgit/semver"
	"autogit/settings/logus"
	"autogit/utils"
)

type SemanticGit struct {
	git *git.Repository
}

func (g *SemanticGit) NewRepo(gitRepo *git.Repository) *SemanticGit {
	g.git = gitRepo
	return g
}

func (g *SemanticGit) GetCurrentVersion() *semver.SemVer {
	returned_vers := &semver.SemVer{
		Major: 0, Minor: 0, Patch: 0,
		AugmentedSemver: semver.AugmentedSemver{Alpha: 0, Beta: 0, Rc: 0},
	}

	latest_hash := g.git.GetLatestCommitHash()
	g.git.ForeachTag(func(tag git.Tag) bool {
		vers, err := semver.Parse(tag.Name)
		utils.CheckWarn(err, "WARN failed to parse tag=", tag.Name)

		if tag.Hash == latest_hash || (vers.Prerelease == "" && tag.Hash == latest_hash) {
			return false
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
			return true
		}

		return false
	})

	return returned_vers
}

func (g *SemanticGit) CalculateNextVersion(vers *semver.SemVer) *semver.SemVer {

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

func (g *SemanticGit) GetNextVersion(semver_options semver.OptionsSemVer) *semver.SemVer {
	vers := g.GetCurrentVersion()
	vers.Options = semver_options
	vers = g.CalculateNextVersion(vers)

	return vers
}

func (g *SemanticGit) GetChangelogByTag(fromTag string, enable_warnings bool) []conventionalcommits.ConventionalCommit {
	var result []conventionalcommits.ConventionalCommit

	g.git.GetLogsFromTag(fromTag, func(log_record git.Log) bool {

		parsed_commit, err := conventionalcommits.ParseCommit(log_record.Msg)
		if err != nil {
			if enable_warnings {
				logus.Warn("unable to parse commit with hash=", logus.CommitHash(log_record.Hash))
			}
			return false
		}

		// attempt to convert to Semver
		var foundSemver *semver.SemVer
		var foundTag git.Tag
		g.git.ForeachTag(func(tag git.Tag) bool {
			if log_record.Hash == tag.Hash {
				semver, err := semver.Parse(tag.Name)
				if err != nil {
					return false
				}
				foundSemver = semver
				foundTag = tag
				return true
			}
			return false
		})
		if foundSemver != nil {
			// Get Changelog only until previous stable tag and don't mind first commit
			if foundSemver.Prerelease == "" && fromTag != foundTag.Name && log_record.Hash != g.git.GetLatestCommitHash() {
				return true
			}
		}

		parsed_commit.Hash = log_record.Hash.String()[:8]
		result = append(result, *parsed_commit)

		return false
	})

	return result
}

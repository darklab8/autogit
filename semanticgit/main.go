/*
Git with applied semantic versioning
and conventional commits to it
*/
package semanticgit

import (
	"autogit/semanticgit/conventionalcommits"
	"autogit/semanticgit/git"
	"autogit/semanticgit/semver"
	"autogit/utils"
	"log"
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

		if tag.Hash == latest_hash && vers.Prerelease == "" {
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

	changelog := g.GetChangelogByTag("", false)

	var major_change, minor_change, patch_change bool
	for _, record := range changelog.Logs {

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

	for _, tag_vers := range changelog.Tags {
		if tag_vers.Alpha > vers.Alpha {
			vers.Alpha = tag_vers.Alpha
		}
		if tag_vers.Beta > vers.Beta {
			vers.Beta = tag_vers.Beta
		}
		if tag_vers.Rc > vers.Rc {
			vers.Rc = tag_vers.Rc
		}
	}

	if vers.Options.Alpha && !vers.Options.Beta && !vers.Options.Rc {
		vers.Alpha++
	}
	if vers.Options.Beta && !vers.Options.Rc {
		vers.Beta++

		if vers.Alpha == 0 {
			vers.Alpha = 1
		}
	}
	if vers.Options.Rc {
		vers.Rc++

		if vers.Alpha == 0 {
			vers.Alpha = 1
		}
		if vers.Beta == 0 {
			vers.Beta = 1
		}
	}

	return vers
}

func (g *SemanticGit) GetNextVersion(semver_options semver.OptionsSemVer) *semver.SemVer {
	vers := g.GetCurrentVersion()
	vers.Options = semver_options
	vers = g.CalculateNextVersion(vers)

	return vers
}

type ChangelogByTagResult struct {
	Logs []conventionalcommits.ConventionalCommit
	Tags []semver.SemVer
}

func (g *SemanticGit) GetChangelogByTag(fromTag string, enable_warnings bool) ChangelogByTagResult {
	result := ChangelogByTagResult{}

	g.git.GetLogsFromTag(fromTag, func(log_record git.Log) bool {

		parsed_commit, err := conventionalcommits.ParseCommit(log_record.Msg)
		if err != nil {
			if enable_warnings {
				log.Println("WARN unable to parse commit with hash=", log_record.Hash.String())
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
			if fromTag != foundTag.Name && log_record.Hash != g.git.GetLatestCommitHash() {
				result.Tags = append(result.Tags, *foundSemver)
			}
			// Get Changelog only until previous stable tag and don't mind first commit
			if foundSemver.Prerelease == "" && fromTag != foundTag.Name && log_record.Hash != g.git.GetLatestCommitHash() {
				return true
			}
		}

		parsed_commit.Hash = log_record.Hash.String()[:8]
		result.Logs = append(result.Logs, *parsed_commit)

		return false
	})

	return result
}

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
	latest_tag := g.git.GetLatestTagString(true)
	if latest_tag == "" {
		return &semver.SemVer{Major: 0, Minor: 0, Patch: 0}
	}

	vers, err := semver.Parse(latest_tag)

	if err != nil {
		utils.CheckFatal(err, "ERR failed to parse latest_tag={%s}\nAutofixing to semantic version being default", latest_tag)
		return &semver.SemVer{}
	}

	return vers
}

func (g *SemanticGit) CalculateNextVersion(vers *semver.SemVer) *semver.SemVer {

	logs := g.GetChangelogByTag("", false)

	var major_change, minor_change, patch_change bool
	for _, record := range logs {

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
		return vers
	}
	if minor_change {
		vers.Minor += 1
		vers.Patch = 0
		return vers
	}
	if patch_change {
		vers.Patch += 1
		return vers
	}

	return vers
}

func (g *SemanticGit) GetNextVersion() *semver.SemVer {
	vers := g.GetCurrentVersion()
	vers = g.CalculateNextVersion(vers)

	return vers
}

func (g *SemanticGit) GetChangelogByTag(tag string, enable_warnings bool) []conventionalcommits.ConventionalCommit {
	logs := g.git.TestGetChangelogByTag(tag)

	var results []conventionalcommits.ConventionalCommit

	for _, log_record := range logs {
		parsed_commit, err := conventionalcommits.ParseCommit(log_record.Msg)
		if err != nil {
			if enable_warnings {
				log.Println("WARN unable to parse commit with hash=", log_record.Hash.String())
			}
			continue
		}
		parsed_commit.Hash = log_record.Hash.String()[:8]
		results = append(results, *parsed_commit)
	}

	return results
}

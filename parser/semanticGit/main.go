/*
Git with applied semantic versioning
and conventional commits to it
*/
package semanticgit

import (
	"autogit/git"
	"autogit/parser/conventionalcommits"
	"autogit/parser/semver"
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
	latest_tag := g.git.GetLatestTagString()
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
	// Calculate next version from changelog additions
	return vers
}

func (g *SemanticGit) GetNextVersion() *semver.SemVer {
	vers := g.GetCurrentVersion()

	vers = g.CalculateNextVersion(vers)

	return vers
}

func (g *SemanticGit) GetChangelogByTag(tag string) []conventionalcommits.ConventionalCommit {
	logs := g.git.TestGetChangelogByTag(tag)

	var results []conventionalcommits.ConventionalCommit

	for _, log_record := range logs {
		parsed_commit, err := conventionalcommits.ParseCommit(log_record.Msg)
		if err != nil {
			log.Println("WARN unable to parse commit with hash={%s}", log_record.Hash.String())
			continue
		}
		results = append(results, *parsed_commit)
	}

	return results
}

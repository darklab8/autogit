package semver

import (
	"log"
	"strconv"

	"github.com/darklab8/autogit/semanticgit/semver/semvertype"
	"github.com/darklab8/autogit/settings"
	"github.com/darklab8/autogit/settings/types"
)

func ParseToInt(msg string) int {
	number, err := strconv.Atoi(msg)
	if err != nil {
		log.Fatal("failed to parse to int", msg)
	}
	return number
}

func Parse(msg types.TagName) (*semvertype.SemVer, error) {
	matched := settings.RegexSemVer.FindStringSubmatch(string(msg))

	if len(matched) == 0 {
		return nil, NotParsedSemver{}
	}

	// Allowing not defining patch always
	patch_version := matched[3]
	if patch_version == "" {
		patch_version = "0"
	}
	vers := &semvertype.SemVer{
		Major:      ParseToInt(matched[1]),
		Minor:      ParseToInt(matched[2]),
		Patch:      ParseToInt(patch_version),
		Prerelease: matched[4],
		Build:      matched[5],
	}

	match_prerelease := settings.RegexPrerelease.FindStringSubmatch(vers.Prerelease)
	if len(match_prerelease) > 0 {

		alpha, alpha_err := strconv.Atoi(match_prerelease[1])
		beta, beta_err := strconv.Atoi(match_prerelease[2])
		rc, rc_err := strconv.Atoi(match_prerelease[3])

		if alpha_err != nil {
			alpha = 0
		}
		if beta_err != nil {
			beta = 0
		}
		if rc_err != nil {
			rc = 0
		}

		vers.Alpha = alpha
		vers.Beta = beta
		vers.Rc = rc
	}

	return vers, nil
}

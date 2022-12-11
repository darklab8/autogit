package semver

import (
	"autogit/utils"
	"regexp"
	"strconv"
)

var semVerRegex *regexp.Regexp

func init() {
	// copied from https://semver.org/spec/v2.0.0.html
	utils.InitRegexExpression(&semVerRegex,
		`^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
}

type SemVer struct {
	// Official SemVer
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
}

func ParseToInt(msg string) int {
	number, err := strconv.Atoi(msg)
	utils.CheckPanic(err)
	return number
}

func Parse(msg string) SemVer {
	matched := semVerRegex.FindStringSubmatch(msg)
	_ = matched
	return SemVer{
		Major:      ParseToInt(matched[1]),
		Minor:      ParseToInt(matched[2]),
		Patch:      ParseToInt(matched[3]),
		Prerelease: matched[4],
		Build:      matched[5],
	}
}

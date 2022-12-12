package semver

import (
	"autogit/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

func Parse(msg string) (*SemVer, error) {
	matched := semVerRegex.FindStringSubmatch(msg)

	if len(matched) == 0 {
		return nil, NotParsed{}
	}

	return &SemVer{
		Major:      ParseToInt(matched[1]),
		Minor:      ParseToInt(matched[2]),
		Patch:      ParseToInt(matched[3]),
		Prerelease: matched[4],
		Build:      matched[5],
	}, nil
}

func (s SemVer) ToString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("v%d", s.Major))
	sb.WriteString(fmt.Sprintf(".%d", s.Minor))
	sb.WriteString(fmt.Sprintf(".%d", s.Patch))

	if s.Prerelease != "" {
		sb.WriteString(fmt.Sprintf("-%s", s.Prerelease))
	}

	if s.Build != "" {
		sb.WriteString(fmt.Sprintf("+%s", s.Build))
	}

	return sb.String()
}
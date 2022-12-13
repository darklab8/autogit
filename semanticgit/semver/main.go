package semver

import (
	"autogit/settings"
	"autogit/utils"
	"fmt"
	"strconv"
	"strings"
)

type SemVer struct {
	DisableVFlag bool
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
	matched := settings.RegexSemVer.FindStringSubmatch(msg)

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
	if !s.DisableVFlag {
		sb.WriteString("v")
	}
	sb.WriteString(fmt.Sprintf("%d", s.Major))

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

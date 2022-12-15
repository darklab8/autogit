package semver

import (
	"autogit/settings"
	"autogit/utils"
	"fmt"
	"strconv"
	"strings"
)

type SemVer struct {
	AugmentedSemver
	// Official SemVer
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
}

type OptionsSemVer struct {
	DisableVFlag  bool
	EnableNewline bool
	Build         string
	Alpha         bool
	Beta          bool
	Rc            bool
}
type AugmentedSemver struct {
	Options OptionsSemVer
	Alpha   int
	Beta    int
	Rc      int
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

	vers := &SemVer{
		Major:      ParseToInt(matched[1]),
		Minor:      ParseToInt(matched[2]),
		Patch:      ParseToInt(matched[3]),
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

func (s SemVer) ToString() string {
	var sb strings.Builder
	if !s.Options.DisableVFlag {
		sb.WriteString("v")
	}
	sb.WriteString(fmt.Sprintf("%d", s.Major))

	sb.WriteString(fmt.Sprintf(".%d", s.Minor))
	sb.WriteString(fmt.Sprintf(".%d", s.Patch))

	if s.Options.Alpha || s.Options.Beta || s.Options.Rc {
		sb.WriteString("-")

		if s.Options.Alpha {
			sb.WriteString(fmt.Sprintf("a.%d", s.Alpha))
		}
		if s.Options.Alpha && s.Options.Beta {
			sb.WriteString(".")
		}
		if s.Options.Beta {
			sb.WriteString(fmt.Sprintf("b.%d", s.Beta))
		}
		if s.Options.Beta && s.Options.Rc {
			sb.WriteString(".")
		}
		if s.Options.Rc {
			sb.WriteString(fmt.Sprintf("rc.%d", s.Rc))
		}
	}

	if s.Build != "" {
		sb.WriteString(fmt.Sprintf("+%s", s.Build))
	}

	if s.Options.EnableNewline {
		sb.WriteString("\n")
	}
	return sb.String()
}

package semvertype

import (
	"fmt"
	"strings"

	"github.com/darklab8/autogit/settings/types"
)

type OptionsSemVer struct {
	DisableVFlag  bool
	EnableNewline bool
	Build         string
	Alpha         bool
	Beta          bool
	Rc            bool
	Publish       bool
}

type AugmentedSemver struct {
	Options OptionsSemVer
	Alpha   int
	Beta    int
	Rc      int
}

type SemVer struct {
	AugmentedSemver
	// Official SemVer
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
}

func (s SemVer) ToString() types.TagName {
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
	return types.TagName(sb.String())
}

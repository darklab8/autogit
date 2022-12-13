package settings

import (
	_ "embed"
)

//go:embed version.txt
var Version string

func init() {
	ChangelogInit()
	RegexInit()
}

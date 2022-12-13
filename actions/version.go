package actions

import (
	"autogit/git"
	sGit "autogit/parser/semanticGit"
	"fmt"
)

var VersionDisableVFlag *bool
var VersionBuildMeta *string

func Version() string {
	g := (&sGit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())
	vers := g.GetNextVersion()
	vers.DisableVFlag = (*VersionDisableVFlag)
	vers.Build = *VersionBuildMeta
	return fmt.Sprintf("%s", vers.ToString())
}

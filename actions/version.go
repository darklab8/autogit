package actions

import (
	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"fmt"
)

var VersionDisableVFlag *bool
var VersionBuildMeta *string

func Version() string {
	g := (&semanticgit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())
	vers := g.GetNextVersion()
	vers.DisableVFlag = (*VersionDisableVFlag)
	vers.Build = *VersionBuildMeta
	return fmt.Sprintf("%s", vers.ToString())
}

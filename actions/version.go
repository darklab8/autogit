package actions

import (
	"autogit/git"
	sGit "autogit/parser/semanticGit"
	"fmt"
)

var VersionDisableVFlag *bool

func Version() string {
	g := (&sGit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())
	vers := g.GetNextVersion()
	vers.DisableVFlag = (*VersionDisableVFlag)
	return fmt.Sprintf("%s\n", vers.ToString())
}

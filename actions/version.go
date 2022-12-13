package actions

import (
	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"strings"
)

var VersionDisableVFlag *bool
var VersionDiableNewLine *bool
var VersionBuildMeta *string

func Version() string {
	g := (&semanticgit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())
	vers := g.GetNextVersion()
	vers.DisableVFlag = (*VersionDisableVFlag)
	vers.Build = *VersionBuildMeta

	var sb strings.Builder
	sb.WriteString(vers.ToString())
	if !*VersionDiableNewLine {
		sb.WriteString("\n")
	}
	return sb.String()
}

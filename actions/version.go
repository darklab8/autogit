package actions

import (
	"autogit/git"
	sGit "autogit/parser/semanticGit"
	"fmt"
)

func Version() string {
	g := (&sGit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())
	return fmt.Sprintf("%s\n", g.GetNextVersion().ToString())
}

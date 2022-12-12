package actions

import (
	"autogit/git"
	sGit "autogit/parser/semanticGit"
	"fmt"
	"strings"
)

func Changelog() string {
	g := (&sGit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())

	logs := g.GetChangelogByTag("")

	var sb strings.Builder

	for _, record := range logs {
		sb.WriteString(fmt.Sprintf("%s(%s): %s\n", record.Type, record.Scope, record.Subject))
	}

	return sb.String()
}

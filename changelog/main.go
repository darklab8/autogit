package changelog

// git log --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit

// git log --pretty=format:'%s' --abbrev-commit
import (
	"autogit/git"
)

func Run() {
	(&git.Repository{}).GetRepoInWorkDir().GetLogs()
}

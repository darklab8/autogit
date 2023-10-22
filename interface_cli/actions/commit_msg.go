package actions

import (
	"autogit/interface_cli/actions/validation"
	"autogit/semanticgit/conventionalcommits"
	"autogit/settings"
	"autogit/settings/logus"
	"autogit/settings/types"
	"autogit/settings/utils"
	"os"
	"path/filepath"
)

func CommmitMsg(args []string) {
	if !settings.GetConfig().Validation.Sections.Hook.CommitMsg.Enabled {
		return
	}

	conf := settings.GetConfig()

	inputFile := types.FilePath(filepath.Join(string(utils.GetProjectDir()), args[0]))
	logus.Debug("Received CommitFile", logus.FilePath(inputFile))

	file, err := os.ReadFile(string(inputFile))
	logus.CheckFatal(err, "Could not read the file due to this error")

	// convert the file binary into a string using string
	fileContent := types.CommitOriginalMsg(string(file))
	logus.Debug("acquired file_content", logus.CommitMessage(fileContent))

	commit, err := conventionalcommits.NewCommit(fileContent)
	logus.CheckError(err, "unable to parse commit to conventional commits standard")

	err = validation.Validate(*commit, conf)
	logus.CheckError(err, "failed to validate commits", logus.Commit(commit.ParsedCommit))

	logus.Info("parsed commit", logus.Commit(commit.ParsedCommit))
}

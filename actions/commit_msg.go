package actions

import (
	"autogit/actions/validation"
	"autogit/semanticgit/conventionalcommits"
	"autogit/settings"
	"autogit/settings/logus"
	"autogit/settings/types"
	"os"
)

func CommmitMsg(args []string) {
	conf := settings.LoadSettings(settings.GetSettingsPath())
	settings.LoadSettings(settings.GetSettingsPath())

	inputFile := types.FilePath(args[0])
	logus.Debug("Received CommitFile", logus.FilePath(inputFile))

	file, err := os.ReadFile(string(inputFile))
	logus.CheckFatal(err, "Could not read the file due to this error")

	// convert the file binary into a string using string
	fileContent := types.CommitMessage(string(file))
	logus.Debug("acquired file_content", logus.CommitMessage(fileContent))

	commit, err := conventionalcommits.NewCommit(fileContent)
	logus.CheckFatal(err, "unable to parse commit to conventional commits standard")

	if settings.GetConfig().Validation.Sections.Hook.CommitMsg.Enabled {
		err := validation.Validate(commit, conf)
		logus.CheckFatal(err, "failed to validate commits")
	}

	logus.Info("parsed commit", logus.Commit(commit.ParsedCommit))
}

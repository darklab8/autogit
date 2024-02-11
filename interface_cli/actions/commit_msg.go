package actions

import (
	"autogit/interface_cli/actions/validation"
	"autogit/semanticgit/conventionalcommits"
	"autogit/settings"
	"autogit/settings/logus"
	"autogit/settings/types"
	"os"
	"path/filepath"

	"github.com/darklab8/go-utils/goutils/utils"
	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

func CommmitMsg(args []string) {
	if !settings.GetConfig().Validation.Sections.Hook.CommitMsg.Enabled {
		return
	}

	conf := settings.GetConfig()

	inputFile := utils_types.FilePath(filepath.Join(string(utils.GetProjectDir()), args[0]))
	logus.Log.Debug("Received CommitFile", utils_logus.FilePath(inputFile))

	file, err := os.ReadFile(string(inputFile))
	logus.Log.CheckFatal(err, "Could not read the file due to this error")

	// convert the file binary into a string using string
	fileContent := types.CommitOriginalMsg(string(file))
	logus.Log.Debug("acquired file_content", logus.CommitMessage(fileContent))

	commit, err := conventionalcommits.NewCommit(fileContent)
	logus.Log.CheckError(err, "unable to parse commit to conventional commits standard")

	err = validation.Validate(*commit, conf)
	logus.Log.CheckError(err, "failed to validate commits", logus.Commit(commit.ParsedCommit))

	logus.Log.Info("parsed commit", logus.Commit(commit.ParsedCommit))
}

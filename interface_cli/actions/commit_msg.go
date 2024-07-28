package actions

import (
	"os"
	"path/filepath"

	"github.com/darklab8/autogit/v2/interface_cli/actions/validation"
	"github.com/darklab8/autogit/v2/semanticgit/conventionalcommits"
	"github.com/darklab8/autogit/v2/settings"
	"github.com/darklab8/autogit/v2/settings/logus"
	"github.com/darklab8/autogit/v2/settings/types"

	"github.com/darklab8/go-utils/utils/utils_logus"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func CommmitMsg(args []string) {
	if !settings.GetConfig().Validation.Sections.Hook.CommitMsg.Enabled {
		return
	}

	conf := settings.GetConfig()

	inputFile := utils_types.FilePath(filepath.Join(string(utils_os.GetProjectDir()), args[0]))
	logus.Log.Debug("Received CommitFile", utils_logus.FilePath(inputFile))

	file, err := os.ReadFile(string(inputFile))
	logus.Log.CheckFatal(err, "Could not read the file due to this error")

	// convert the file binary into a string using string
	fileContent := types.CommitOriginalMsg(string(file))
	logus.Log.Debug("acquired file_content", logus.CommitMessage(fileContent))

	commit, err := conventionalcommits.NewCommit(fileContent)
	logus.Log.CheckFatal(err, "unable to parse commit to conventional commits standard")

	err = validation.Validate(*commit, conf)
	logus.Log.CheckFatal(err, "failed to validate commits", logus.Commit(commit.ParsedCommit))

	logus.Log.Info("parsed commit", logus.Commit(commit.ParsedCommit))
}

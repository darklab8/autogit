package actions

import (
	"autogit/actions/validation"
	"autogit/semanticgit/conventionalcommits"
	"autogit/settings"
	"autogit/settings/logus"
	"fmt"
	"os"
)

func CommmitMsg(args []string) {
	conf := settings.LoadSettings(settings.GetSettingsPath())
	settings.LoadSettings(settings.GetSettingsPath())

	inputFile := args[0]
	fmt.Printf("commitFile=%s\n", inputFile)

	file, err := os.ReadFile(inputFile)
	logus.CheckFatal(err, "Could not read the file due to this error")

	// convert the file binary into a string using string
	fileContent := string(file)
	fmt.Printf("fileContent=%s", fileContent)

	commit, err := conventionalcommits.NewCommit(fileContent)
	logus.CheckFatal(err, "unable to parse commit to conventional commits standard")

	if settings.GetConfig().Validation.Sections.Hook.CommitMsg.Enabled {
		err := validation.Validate(commit, conf)
		logus.CheckFatal(err, "failed to validate commits")
	}

	fmt.Printf("parsed commit:\n%s\n", commit.StringAnnotated())
}

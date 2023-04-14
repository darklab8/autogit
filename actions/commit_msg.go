package actions

import (
	"autogit/actions/validation"
	"autogit/semanticgit/conventionalcommits"
	"autogit/settings"
	"autogit/utils"
	"fmt"
	"io/ioutil"
)

func CommmitMsg(args []string) {
	conf := settings.LoadSettings(settings.GetSettingsPath())
	settings.LoadSettings(settings.GetSettingsPath())

	inputFile := args[0]
	fmt.Printf("commitFile=%s\n", inputFile)

	file, err := ioutil.ReadFile(inputFile)
	utils.CheckFatal(err, "Could not read the file due to this %s error \n")

	// convert the file binary into a string using string
	fileContent := string(file)
	fmt.Printf("fileContent=%s", fileContent)

	commit, err := conventionalcommits.NewCommit(fileContent)
	utils.CheckFatal(err, "unable to parse commit to conventional commits standard")

	if settings.GetConfig().Validation.Sections.Hook.CommitMsg.Enabled {
		err := validation.Validate(commit, conf)
		utils.CheckFatal(err)
	}

	fmt.Printf("parsed commit:\n%s\n", commit.StringAnnotated())
}

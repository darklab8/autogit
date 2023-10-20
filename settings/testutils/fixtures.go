package testutils

import (
	"autogit/settings"
	"autogit/settings/types"
	"autogit/settings/utils"
	"path/filepath"
)

func FixtureSettings() {
	workdir := utils.GetCurrrentFolder()
	originalSettingsPath := workdir
	rootFolder := filepath.Dir(string(originalSettingsPath))
	testSettingsPath := types.ConfigPath(filepath.Join(rootFolder, "settings", "autogit.example.yml"))

	config := settings.LoadSettings(testSettingsPath)
	settings.RegexInit(config)
}

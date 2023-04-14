package settings

import (
	"autogit/utils"
	"path/filepath"
)

func FixtureSettings() {
	workdir := utils.GetCurrrentFolder()
	originalSettingsPath := workdir
	rootFolder := filepath.Dir(string(originalSettingsPath))
	testSettingsPath := SettingPath(filepath.Join(rootFolder, "settings", "autogit.example.yml"))

	config := LoadSettings(testSettingsPath)
	RegexInit(config)
}

func FixtureAutouse() {
	FixtureSettings()
}

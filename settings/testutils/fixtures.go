package testutils

import (
	"autogit/settings"
	"autogit/settings/types"
	"path/filepath"

	"github.com/darklab8/go-utils/goutils/utils"
)

func FixtureSettings() {
	workdir := utils.GetCurrentFolder()
	originalSettingsPath := workdir
	rootFolder := filepath.Dir(string(originalSettingsPath))
	testSettingsPath := types.ConfigPath(filepath.Join(rootFolder, "settings", "autogit.example.yml"))

	settings.NewConfig(testSettingsPath)
}

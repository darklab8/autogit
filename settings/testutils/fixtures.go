package testutils

import (
	"path/filepath"

	"github.com/darklab8/autogit/v2/settings"
	"github.com/darklab8/autogit/v2/settings/types"

	"github.com/darklab8/go-utils/goutils/utils"
)

func FixtureSettings() {
	workdir := utils.GetCurrentFolder()
	originalSettingsPath := workdir
	rootFolder := filepath.Dir(string(originalSettingsPath))
	testSettingsPath := types.ConfigPath(filepath.Join(rootFolder, "settings", "autogit.example.yml"))

	settings.NewConfig(testSettingsPath)
}

package testutils

import (
	"path/filepath"

	"github.com/darklab8/autogit/v2/settings"
	"github.com/darklab8/autogit/v2/settings/types"
	"github.com/darklab8/go-utils/utils/utils_os"
)

func FixtureSettings() {
	workdir := utils_os.GetCurrentFolder()
	originalSettingsPath := workdir
	rootFolder := filepath.Dir(string(originalSettingsPath))
	testSettingsPath := types.ConfigPath(filepath.Join(rootFolder, "settings", "autogit.example.yml"))

	settings.NewConfig(testSettingsPath)
}

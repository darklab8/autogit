package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/darklab8/autogit/settings/logus"
	"github.com/darklab8/autogit/settings/types"

	"github.com/darklab8/go-utils/goutils/utils"
)

func TestWriteCommentedConfig(t *testing.T) {
	current_folder := utils.GetCurrentFolder()
	temp_data := filepath.Join(string(current_folder), "temp_data")
	err := os.MkdirAll(temp_data, os.ModePerm)
	logus.Log.CheckFatal(err, "failed to create temp_data folder")

	temp_config := filepath.Join(temp_data, fmt.Sprintf("%s.yml", utils.TokenHex(8)))
	defer os.Remove(temp_config)
	init_write_config(types.ConfigPath(temp_config))
}

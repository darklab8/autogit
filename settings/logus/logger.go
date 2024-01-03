package logus

import (
	"autogit/settings/envs"

	"github.com/darklab8/darklab_goutils/goutils/logus_core"
	"github.com/darklab8/darklab_goutils/goutils/logus_core/logus_types"
)

var (
	Log *logus_core.Logger
)

func init() {

	Log = logus_core.NewLogger(
		logus_types.LogLevel(envs.LogLevel),
		logus_types.EnableJsonFormat(envs.LogTurnJSONLogging),
		logus_types.EnableFileShowing(envs.LogShowFileLocations))
}

package interface_cli

import (
	"fmt"

	"github.com/darklab8/autogit/settings/logus"

	"github.com/darklab8/go-typelog/typelog"
	"github.com/spf13/cobra"
)

type sharedFlags struct {
	verboseLogging *bool
}

func (s *sharedFlags) Bind(Cmd *cobra.Command) {
	s.verboseLogging = Cmd.Flags().BoolP("verbose", "v", false, "Turn on verbose logging")
}

func (s *sharedFlags) Run() {
	if *(s.verboseLogging) {
		logus.Log = typelog.NewLogger("github.com/darklab8/autogit", typelog.WithLogLevel(typelog.LEVEL_DEBUG))
	}
	logus.Log.Debug(fmt.Sprintf("verbose=%t\n", *(s.verboseLogging)))
}

var shared struct {
	version         sharedFlags
	hook_activate   sharedFlags
	hook_deactivate sharedFlags
	hook_commit_msg sharedFlags
	changelog       sharedFlags
	init            sharedFlags
	about           sharedFlags
}

func init() {
	shared.version.Bind(semverCmd)
	shared.hook_activate.Bind(activateCmd)
	shared.hook_deactivate.Bind(deactivateCmd)
	shared.hook_commit_msg.Bind(commitMsgCmd)
	shared.changelog.Bind(changelogCmd)
	shared.init.Bind(initCmd)
	shared.about.Bind(aboutCmd)
}

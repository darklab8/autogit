package actions

import (
	"autogit/settings"
	"fmt"
)

func Version() string {
	return fmt.Sprintf("version: v%s\n", settings.Version)
}

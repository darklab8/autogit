package actions

import (
	"autogit/settings"
	"fmt"
	"strings"
)

func About() string {
	var sb strings.Builder
	sb.WriteString("OK ")
	sb.WriteString("autogit version: ")
	sb.WriteString(fmt.Sprintf("%s", settings.Version))
	return sb.String()
}

package actions

import (
	"autogit/settings"
	"strings"
)

func About() string {
	var sb strings.Builder
	sb.WriteString("OK ")
	sb.WriteString("autogit version: ")
	sb.WriteString(settings.Version)
	return sb.String()
}

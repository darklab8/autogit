package utils

import (
	"autogit/settings/logus"
	"autogit/settings/types"
	"regexp"
)

func InitRegexExpression(regex **regexp.Regexp, expression types.RegexExpression) {
	var err error

	*regex, err = regexp.Compile(string(expression))
	logus.CheckFatal(err, "failed to init regex={%s} in ", logus.Regex(expression), logus.FilePath(GetCurrentFile()))
}

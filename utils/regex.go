package utils

import (
	"regexp"
)

func InitRegexExpression(regex **regexp.Regexp, expression string) {
	var err error

	*regex, err = regexp.Compile(expression)
	CheckFatal(err, "failed to parse numberParser in ", GetCurrentFile())
}

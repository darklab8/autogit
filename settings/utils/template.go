package utils

import (
	"autogit/settings/logus"
	"autogit/settings/types"
	"bytes"
	"text/template"
)

func TmpRender(templateRef *template.Template, data interface{}) string {
	var header bytes.Buffer
	err := templateRef.Execute(&header, data)
	logus.CheckFatal(err, "failed executing regex")
	return header.String()
}

func TmpInit(content types.TemplateExpression) *template.Template {
	var err error
	templateRef, err := template.New("test").Parse(string(content))
	logus.CheckFatal(err, "failed initing regex template")
	return templateRef
}

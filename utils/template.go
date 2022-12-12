package utils

import (
	"bytes"
	"text/template"
)

func TmpRender(templateRef *template.Template, data interface{}) string {
	var header bytes.Buffer
	err := templateRef.Execute(&header, data)
	CheckFatal(err)
	return header.String()
}

func TmpInit(content string) *template.Template {
	var err error
	templateRef, err := template.New("test").Parse(content)
	CheckFatal(err)
	return templateRef
}

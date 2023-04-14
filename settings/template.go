package settings

import (
	"autogit/utils"
	_ "embed"
	"text/template"
)

type Templates struct {
	CommitUrl      *template.Template
	CommitRangeUrl *template.Template
	IssueUrl       *template.Template
}

func initTemplates(conf ConfigScheme) *Templates {
	result := &Templates{}
	result.CommitUrl = utils.TmpInit(conf.Changelog.CommitURL)
	result.CommitRangeUrl = utils.TmpInit(conf.Changelog.CommitRangeURL)
	result.IssueUrl = utils.TmpInit(conf.Changelog.IssueURL)
	return result
}

var templates *Templates

func GetTemplates() Templates {
	conf := GetConfig()
	if templates == nil {
		templates = initTemplates(conf)
	}
	return *templates
}

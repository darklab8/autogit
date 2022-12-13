package settings

import (
	"autogit/utils"
	_ "embed"
	"text/template"
)

var Template struct {
	CommitUrl      *template.Template
	CommitRangeUrl *template.Template
	IssueUrl       *template.Template
}

func init() {
	Template.CommitUrl = utils.TmpInit(Config.Changelog.CommitURL)
	Template.CommitRangeUrl = utils.TmpInit(Config.Changelog.CommitRangeURL)
	Template.IssueUrl = utils.TmpInit(Config.Changelog.IssueURL)
}

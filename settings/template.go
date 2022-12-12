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
	Changelog      *template.Template
}

//go:embed templates/changelog.md
var changelogMarkup string
var changelogTemplate *template.Template

func init() {
	Template.Changelog = utils.TmpInit(changelogMarkup)

	Template.CommitUrl = utils.TmpInit(Config.Changelog.CommitURL)
	Template.CommitRangeUrl = utils.TmpInit(Config.Changelog.CommitRangeURL)
	Template.IssueUrl = utils.TmpInit(Config.Changelog.IssueURL)
}

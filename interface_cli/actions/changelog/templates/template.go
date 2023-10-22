package templates

import (
	"autogit/semanticgit/conventionalcommits"
	"autogit/settings"
	"autogit/settings/types"
	"autogit/settings/utils"
	_ "embed"
	"fmt"
	"text/template"
	"time"
)

type Templates struct {
	commitUrl struct {
		Template *template.Template
	}
	commitRangeUrl struct {
		Template *template.Template
	}
	issueUrl struct {
		Template *template.Template
	}
	conf settings.ChangelogScheme
}

func initTemplates(conf settings.ChangelogScheme) *Templates {
	result := &Templates{conf: conf}
	result.commitUrl.Template = utils.TmpInit(conf.CommitURL)
	result.commitRangeUrl.Template = utils.TmpInit(conf.CommitRangeURL)
	result.issueUrl.Template = utils.TmpInit(conf.IssueURL)
	return result
}

var cached_templates *Templates

func NewTemplates() Templates {
	conf := settings.GetConfig()
	if cached_templates == nil {
		cached_templates = initTemplates(conf.Changelog)
	}
	return *cached_templates
}

type CommitUrlVars struct {
	Commit           string
	REPOSITORY_OWNER string
	REPOSITORY_NAME  string
}

func (templs Templates) RenderCommitUrl(record conventionalcommits.ConventionalCommit) string {
	return utils.TmpRender(templs.commitUrl.Template, CommitUrlVars{
		Commit:           record.Hash,
		REPOSITORY_OWNER: templs.conf.REPOSITORY_OWNER,
		REPOSITORY_NAME:  templs.conf.REPOSITORY_NAME,
	})
}

type IssueDataVars struct {
	Issue            string
	REPOSITORY_OWNER string
	REPOSITORY_NAME  string
}

func (templs Templates) RenderIssueUrl(issue_n string) string {
	return utils.TmpRender(
		templs.issueUrl.Template,
		IssueDataVars{Issue: issue_n,
			REPOSITORY_OWNER: templs.conf.REPOSITORY_OWNER,
			REPOSITORY_NAME:  templs.conf.REPOSITORY_NAME},
	)
}

type CommitRangeUrlVars struct {
	From             string
	To               string
	REPOSITORY_OWNER string
	REPOSITORY_NAME  string
}

func (templs Templates) NewCommitRangeUrlRender(logs []conventionalcommits.ConventionalCommit, ChangelogVersion types.TagName) string {
	r := CommitRangeUrlVars{
		From:             logs[len(logs)-1].Hash,
		To:               logs[0].Hash,
		REPOSITORY_OWNER: templs.conf.REPOSITORY_OWNER,
		REPOSITORY_NAME:  templs.conf.REPOSITORY_NAME,
	}
	currentTime := time.Now()
	return fmt.Sprintf("## **%s** <sub><sub>%s ([%s...%s](%s))</sub></sub>", ChangelogVersion, currentTime.Format("2006-01-02"), r.From, r.To, utils.TmpRender(templs.commitRangeUrl.Template, r))
}

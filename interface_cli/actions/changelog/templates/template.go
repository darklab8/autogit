package templates

import (
	_ "embed"
	"text/template"
	"time"

	"github.com/darklab8/autogit/v2/semanticgit/conventionalcommits"
	"github.com/darklab8/autogit/v2/semanticgit/conventionalcommits/conventionalcommitstype"
	"github.com/darklab8/autogit/v2/settings"
	"github.com/darklab8/autogit/v2/settings/logus"
	"github.com/darklab8/autogit/v2/settings/types"

	"github.com/darklab8/go-utils/utils"
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
	CommitHash       conventionalcommitstype.Hash
	REPOSITORY_OWNER string
	REPOSITORY_NAME  string
}

func (templs Templates) RenderCommitUrl(record conventionalcommits.ConventionalCommit) string {
	return utils.TmpRender(templs.commitUrl.Template, CommitUrlVars{
		CommitHash:       record.Hash,
		REPOSITORY_OWNER: templs.conf.REPOSITORY_OWNER,
		REPOSITORY_NAME:  templs.conf.REPOSITORY_NAME,
	})
}

type IssueDataVars struct {
	Issue            conventionalcommitstype.Issue
	REPOSITORY_OWNER string
	REPOSITORY_NAME  string
}

func (templs Templates) RenderIssueUrl(issue_n conventionalcommitstype.Issue) string {
	return utils.TmpRender(
		templs.issueUrl.Template,
		IssueDataVars{Issue: issue_n,
			REPOSITORY_OWNER: templs.conf.REPOSITORY_OWNER,
			REPOSITORY_NAME:  templs.conf.REPOSITORY_NAME},
	)
}

type CommitRangeUrlVars struct {
	From             conventionalcommitstype.Hash
	To               conventionalcommitstype.Hash
	REPOSITORY_OWNER string
	REPOSITORY_NAME  string
}

func (templs Templates) NewCommitRangeUrlRender(logs []conventionalcommits.ConventionalCommit, ChangelogVersion types.TagName) Header {
	var from, to conventionalcommitstype.Hash
	if len(logs) == 0 {
		logus.Log.Error("for some reason logs count is 0 at NewCommitRangeUrlRender")
		from = "undefined"
		to = "undefined"
	} else {
		from = logs[len(logs)-1].Hash
		to = logs[0].Hash
	}

	r := CommitRangeUrlVars{
		From:             from,
		To:               to,
		REPOSITORY_OWNER: templs.conf.REPOSITORY_OWNER,
		REPOSITORY_NAME:  templs.conf.REPOSITORY_NAME,
	}
	currentTime := time.Now()
	return Header{
		ChangelogVersion: ChangelogVersion,
		Timestamp:        currentTime.Format("2006-01-02"),
		From:             r.From,
		To:               r.To,
		CommitRangeURL:   utils.TmpRender(templs.commitRangeUrl.Template, r),
	}
}

type Header struct {
	ChangelogVersion types.TagName
	Timestamp        string
	From             conventionalcommitstype.Hash
	To               conventionalcommitstype.Hash
	CommitRangeURL   string
}

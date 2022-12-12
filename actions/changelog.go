package actions

import (
	"autogit/git"
	"autogit/parser/conventionalcommits"
	sGit "autogit/parser/semanticGit"
	"autogit/utils"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
	"time"

	_ "embed"

	"gopkg.in/yaml.v3"
)

type ConfigScheme struct {
	Changelog struct {
		CommitURL      string `yaml:"commitUrl"`
		CommitRangeURL string `yaml:"commitRangeUrl"`
		IssueURL       string `yaml:"issueUrl"`
	} `yaml:"changelog"`
}

func Changelog() string {
	file, err := ioutil.ReadFile("autogit.yml")
	if err != nil {
		fmt.Printf("Could not read the file due to this %s error \n", err)
	}
	config := ConfigScheme{}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("error: ", err)
	}

	g := (&sGit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())

	logs := g.GetChangelogByTag("", true)

	var commitUrl *template.Template = utils.TmpInit(config.Changelog.CommitURL)
	var commitRangeUrl *template.Template = utils.TmpInit(config.Changelog.CommitRangeURL)
	var IssueUrl *template.Template = utils.TmpInit(config.Changelog.IssueURL)

	var Range struct {
		From string
		To   string
	}
	Range.From = logs[len(logs)-1].Hash
	Range.To = logs[0].Hash

	currentTime := time.Now()

	templateData := changelogData{
		Version: fmt.Sprintf("## **%s** <sub><sub>%s ([%s...%s](%s))</sub></sub>", g.GetNextVersion().ToString(), currentTime.Format("2006-01-02"), Range.From, Range.To, utils.TmpRender(commitRangeUrl, Range)),
	}

	type commitRecord struct {
		Commit string
	}

	for _, record := range logs {
		var issue_rendered strings.Builder
		for _, issue_n := range record.Issue {
			issue_rendered.WriteString(fmt.Sprintf(", [#%s](%s)", issue_n, utils.TmpRender(IssueUrl, struct{ Issue string }{Issue: issue_n})))
		}

		rendered_subject := record.Subject
		IssueMatch := conventionalcommits.IssueRegex.FindAllStringSubmatch(record.Subject, -1)
		for _, match := range IssueMatch {
			rendered_subject = strings.Replace(rendered_subject, match[0], fmt.Sprintf("[#%s](%s)", match[1], utils.TmpRender(IssueUrl, struct{ Issue string }{Issue: match[1]})), -1)
		}

		formatted_url := utils.TmpRender(commitUrl, commitRecord{Commit: record.Hash})
		formatted := fmt.Sprintf("* %s ([%s](%s)%s)\n", rendered_subject, record.Hash, formatted_url, issue_rendered.String())
		if record.Type == "feat" {
			templateData.Features = append(templateData.Features, formatted)
		} else if record.Type == "fix" {
			templateData.Fixes = append(templateData.Fixes, formatted)
		}
	}

	return utils.TmpRender(changelogTemplate, templateData)
}

type changelogData struct {
	Version  string
	Features []string
	Fixes    []string
}

//go:embed templates/changelog.md
var changelogMarkup string
var changelogTemplate *template.Template

func init() {
	changelogTemplate = utils.TmpInit(changelogMarkup)
}

package actions

import (
	"autogit/git"
	sGit "autogit/parser/semanticGit"
	"autogit/utils"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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

	var commitUrl *template.Template = initTemplate(config.Changelog.CommitURL)
	var commitRangeUrl *template.Template = initTemplate(config.Changelog.CommitRangeURL)
	var IssueUrl *template.Template = initTemplate(config.Changelog.IssueURL)

	var Range struct {
		From string
		To   string
	}
	Range.From = logs[len(logs)-1].Hash
	Range.To = logs[0].Hash

	currentTime := time.Now()

	templateData := changelogData{
		Version: fmt.Sprintf("## **%s** <sub><sub>%s ([%s...%s](%s))</sub></sub>", g.GetNextVersion().ToString(), currentTime.Format("2006-01-02"), Range.From, Range.To, Render(commitRangeUrl, Range)),
	}

	type commitRecord struct {
		Commit string
	}

	for _, record := range logs {
		issue_rendered := ""
		if record.Issue != "" {
			issue_rendered = fmt.Sprintf(", [#%s](%s)", record.Issue, Render(IssueUrl, struct{ Issue string }{Issue: record.Issue}))
		}
		formatted_url := Render(commitUrl, commitRecord{Commit: record.Hash})
		formatted := fmt.Sprintf("* %s ([%s](%s)%s)\n", record.Subject, record.Hash, formatted_url, issue_rendered)
		if record.Type == "feat" {
			templateData.Features = append(templateData.Features, formatted)
		} else if record.Type == "fix" {
			templateData.Fixes = append(templateData.Fixes, formatted)
		}
	}

	return Render(changelogTemplate, templateData)
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
	changelogTemplate = initTemplate(changelogMarkup)
}

func Render(templateRef *template.Template, data interface{}) string {
	var header bytes.Buffer
	err := templateRef.Execute(&header, data)
	utils.CheckFatal(err)
	return header.String()
}

func initTemplate(content string) *template.Template {
	var err error
	templateRef, err := template.New("test").Parse(content)
	utils.CheckFatal(err)
	return templateRef
}

package actions

import (
	"autogit/git"
	sGit "autogit/parser/semanticGit"
	"autogit/utils"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"

	_ "embed"

	"gopkg.in/yaml.v3"
)

type ConfigScheme struct {
	View struct {
		CommitURL      string `yaml:"commitUrl"`
		CommitRangeURL string `yaml:"commitRangeUrl"`
		IssueURL       string `yaml:"issueUrl"`
	} `yaml:"view"`
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

	templateData := changelogData{
		Version: fmt.Sprintf("## **%s**", g.GetNextVersion().ToString()),
	}

	var commitUrl *template.Template = initTemplate(config.View.CommitURL)
	// var commitRangeUrl *template.Template = initTemplate(config.View.CommitRangeURL)
	// var IssueUrl *template.Template = initTemplate(config.View.IssueURL)

	type commitRecord struct {
		Commit string
	}

	for _, record := range logs {
		if record.Type == "feat" {
			formatted_url := Render(commitUrl, commitRecord{Commit: record.Hash})
			formatted := fmt.Sprintf("* %s ([%s](%s))\n", record.Subject, record.Hash, formatted_url)
			templateData.Features = append(templateData.Features, formatted)
		} else if record.Type == "fix" {
			formatted_url := Render(commitUrl, commitRecord{Commit: record.Hash})
			formatted := fmt.Sprintf("* %s ([%s](%s))\n", record.Subject, record.Hash, formatted_url)
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

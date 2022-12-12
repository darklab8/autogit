package actions

import (
	"autogit/git"
	"autogit/parser/conventionalcommits"
	sGit "autogit/parser/semanticGit"
	"autogit/utils"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"strings"

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
	// fmt.Printf("--- t:\n%v\n\n", config)

	g := (&sGit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())

	logs := g.GetChangelogByTag("", true)

	var features []conventionalcommits.ConventionalCommit
	var fixes []conventionalcommits.ConventionalCommit

	for _, record := range logs {
		if record.Type == "feat" {
			features = append(features, record)
		} else if record.Type == "fix" {
			fixes = append(fixes, record)
		}
	}

	var sb strings.Builder
	sb.WriteString(Render(HeaderTemplate, headerData{Version: g.GetNextVersion().ToString()}))

	// if len(features) > 0 {
	// 	sb.WriteString("### Features")
	// }

	// if len(fixes) > 0 {
	// 	sb.WriteString("### Bug Fixes")
	// }

	return sb.String()
}

type headerData struct {
	Version string
}

//go:embed templates/header.md
var headerMarkup string
var HeaderTemplate *template.Template

func init() {
	initTemplate(&HeaderTemplate, headerMarkup)
}

func Render(templateRef *template.Template, data interface{}) string {
	var header bytes.Buffer
	err := templateRef.Execute(&header, data)
	utils.CheckFatal(err)
	return header.String()
}

func initTemplate(templateRef **template.Template, content string) {
	var err error
	*templateRef, err = template.New("test").Parse(content)
	utils.CheckFatal(err)
}

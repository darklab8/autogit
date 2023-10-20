package logus

import (
	"autogit/semanticgit/conventionalcommits/conventionalcommitstype"
	"autogit/settings/types"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/go-git/go-git/v5/plumbing"
)

func logGroupFiles() slog.Attr {
	return slog.Group("files",
		"file3", GetCallingFile(3),
		"file4", GetCallingFile(4),
	)
}

type slogGroup struct {
	params map[string]string
}

func (s slogGroup) Render() slog.Attr {
	anies := []any{}
	for key, value := range s.params {
		anies = append(anies, key)
		anies = append(anies, value)
	}

	return slog.Group("extras", anies...)
}

type slogParam func(r *slogGroup)

func newSlogGroup(opts ...slogParam) slog.Attr {
	client := &slogGroup{params: make(map[string]string)}
	for _, opt := range opts {
		opt(client)
	}

	return (*client).Render()
}

func TestParam(value int) slogParam {
	return func(c *slogGroup) {
		c.params["test_param"] = fmt.Sprintf("%d", value)
	}
}

func ConfigPath(value types.ConfigPath) slogParam {
	return func(c *slogGroup) {
		c.params["config_path"] = string(value)
	}
}

func Expected(value any) slogParam {
	return func(c *slogGroup) {
		c.params["expected"] = fmt.Sprintf("%v", value)
	}
}
func Actual(value any) slogParam {
	return func(c *slogGroup) {
		c.params["actual"] = fmt.Sprintf("%v", value)
	}
}

func CommitHash(value plumbing.Hash) slogParam {
	return func(c *slogGroup) {
		c.params["commit_hash"] = value.String()
	}
}

func TagName(value types.TagName) slogParam {
	return func(c *slogGroup) {
		c.params["tag_name"] = string(value)
	}
}

func ProjectFolder(value types.ProjectFolder) slogParam {
	return func(c *slogGroup) {
		c.params["project_folder"] = string(value)
	}
}

func FilePath(value types.FilePath) slogParam {
	return func(c *slogGroup) {
		c.params["file_path"] = string(value)
	}
}

func Regex(value types.RegexExpression) slogParam {
	return func(c *slogGroup) {
		c.params["regex"] = string(value)
	}
}

func CommitMessage(value types.CommitMessage) slogParam {
	return func(c *slogGroup) {
		c.params["commit_file"] = string(value)
	}
}

func Commit(commit conventionalcommitstype.ParsedCommit) slogParam {
	return func(c *slogGroup) {
		c.params["commit_type"] = commit.Type
		c.params["commit_scope"] = commit.Scope
		c.params["commit_subject"] = commit.Subject
		c.params["commit_body"] = commit.Body
		c.params["commit_exlamation"] = strconv.FormatBool(commit.Exclamation)
		c.params["commit_hash"] = commit.Hash
		for index, footer := range commit.Footers {
			// Should have made structured logging allowing nested dictionaries.
			// Using as work around more lazy option
			c.params[fmt.Sprintf("commit_footer_%d", index)] = fmt.Sprintf(
				"footer #%d - token: %s, content: %s\n",
				index,
				footer.Token,
				footer.Content,
			)
		}
		for index, issue := range commit.Issue {
			// Should have made structured logging allowing nested dictionaries.
			// Using as work around more lazy option
			c.params[fmt.Sprintf("commit_issue_%d", index)] = issue
		}
	}
}

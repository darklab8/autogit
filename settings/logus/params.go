package logus

import (
	"autogit/settings/types"
	"fmt"
	"log/slog"

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

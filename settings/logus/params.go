package logus

import (
	"autogit/semanticgit/conventionalcommits/conventionalcommitstype"
	"autogit/semanticgit/semver/semvertype"
	"autogit/settings/types"
	"fmt"
	"strconv"

	"github.com/darklab8/darklab_goutils/goutils/logus_core"
	"github.com/go-git/go-git/v5/plumbing"
)

func ConfigPath(value types.ConfigPath) logus_core.SlogParam {
	return func(c *logus_core.SlogGroup) {
		c.Params["config_path"] = string(value)
	}
}

func CommitHash(value plumbing.Hash) logus_core.SlogParam {
	return func(c *logus_core.SlogGroup) {
		c.Params["commit_hash"] = value.String()
	}
}

func TagName(value types.TagName) logus_core.SlogParam {
	return func(c *logus_core.SlogGroup) {
		c.Params["tag_name"] = string(value)
	}
}

func ProjectFolder(value types.ProjectFolder) logus_core.SlogParam {
	return func(c *logus_core.SlogGroup) {
		c.Params["project_folder"] = string(value)
	}
}

func CommitMessage(value types.CommitOriginalMsg) logus_core.SlogParam {
	return func(c *logus_core.SlogGroup) {
		c.Params["commit_file"] = string(value)
	}
}

func SettingsKey(value any) logus_core.SlogParam {
	return func(c *logus_core.SlogGroup) {
		c.Params["settings_key"] = fmt.Sprintf("%v", value)
	}
}

func Commit(commit conventionalcommitstype.ParsedCommit) logus_core.SlogParam {
	return func(c *logus_core.SlogGroup) {
		c.Params["commit_type"] = string(commit.Type)
		c.Params["commit_scope"] = string(commit.Scope)
		c.Params["commit_subject"] = string(commit.Subject)
		c.Params["commit_body"] = string(commit.Body)
		c.Params["commit_exlamation"] = strconv.FormatBool(commit.Exclamation)
		c.Params["commit_hash"] = string(commit.Hash)
		for index, footer := range commit.Footers {
			// Should have made structured logging allowing nested dictionaries.
			// Using as work around more lazy option
			c.Params[fmt.Sprintf("commit_footer_%d", index)] = fmt.Sprintf(
				"token: %s, content: %s",
				footer.Token,
				footer.Content,
			)
		}
		for index, issue := range commit.Issue {
			// Should have made structured logging allowing nested dictionaries.
			// Using as work around more lazy option
			c.Params[fmt.Sprintf("commit_issue_%d", index)] = string(issue)
		}
	}
}

func Semver(semver *semvertype.SemVer) logus_core.SlogParam {
	return func(c *logus_core.SlogGroup) {
		if semver == nil {
			c.Params["semver"] = "nil"
			return
		}
		c.Params["semver"] = string(semver.ToString())
	}
}

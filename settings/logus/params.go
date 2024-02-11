package logus

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/darklab8/autogit/semanticgit/conventionalcommits/conventionalcommitstype"
	"github.com/darklab8/autogit/semanticgit/semver/semvertype"
	"github.com/darklab8/autogit/settings/types"

	"github.com/darklab8/go-typelog/typelog"
	"github.com/go-git/go-git/v5/plumbing"
)

func ConfigPath(value types.ConfigPath) typelog.LogType {
	return typelog.String("config_path", string(value))
}

func CommitHash(value plumbing.Hash) typelog.LogType {
	return typelog.String("commit_hash", value.String())
}

func TagName(value types.TagName) typelog.LogType {
	return typelog.String("tag_name", string(value))
}

func ProjectFolder(value types.ProjectFolder) typelog.LogType {
	return typelog.String("project_folder", string(value))
}

func CommitMessage(value types.CommitOriginalMsg) typelog.LogType {
	return typelog.String("commit_file", string(value))
}

func SettingsKey(value any) typelog.LogType {
	return typelog.Any("settings_key", value)
}

func Commit(commit conventionalcommitstype.ParsedCommit) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(
			slog.String("commit_type", string(commit.Type)),
			slog.String("commit_scope", string(commit.Scope)),
			slog.String("commit_subject", string(commit.Subject)),
			slog.String("commit_body", string(commit.Body)),
			slog.String("commit_exlamation", strconv.FormatBool(commit.Exclamation)),
			slog.String("commit_hash", string(commit.Hash)),
		)

		for index, footer := range commit.Footers {
			// Should have made structured logging allowing nested dictionaries.
			// Using as work around more lazy option
			c.Append(
				slog.String(
					fmt.Sprintf("commit_footer_%d", index),
					fmt.Sprintf(
						"token: %s, content: %s",
						footer.Token,
						footer.Content,
					),
				),
			)
		}
		for index, issue := range commit.Issue {
			// Should have made structured logging allowing nested dictionaries.
			// Using as work around more lazy option
			c.Append(
				slog.String(fmt.Sprintf("commit_issue_%d", index), string(issue)),
			)
		}
	}
}

func Semver(semver *semvertype.SemVer) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		if semver == nil {
			c.Append(slog.String("semver", "nil"))
			return
		}
		c.Append(slog.String("semver", string(semver.ToString())))
	}
}

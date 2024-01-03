package types

import "github.com/darklab8/darklab_goutils/goutils/utils/utils_types"

type ConfigPath utils_types.FilePath

func (c ConfigPath) ToFilePath() utils_types.FilePath {
	return utils_types.FilePath(c)
}

type ProjectFolder utils_types.FilePath

type TagName string

type CommitOriginalMsg string

func (c CommitOriginalMsg) ToString() string {
	return string(c)
}

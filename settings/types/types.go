package types

type FilePath string

type ConfigPath FilePath

func (c ConfigPath) ToFilePath() FilePath {
	return FilePath(c)
}

type ProjectFolder FilePath

type RegexExpression string

type TagName string

type CommitOriginalMsg string

func (c CommitOriginalMsg) ToString() string {
	return string(c)
}

type LogLevel string

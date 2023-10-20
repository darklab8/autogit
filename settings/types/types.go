package types

type ConfigPath FilePath

func (c ConfigPath) ToFilePath() FilePath {
	return FilePath(c)
}

type ProjectFolder string

type FilePath string

type RegexExpression string

type TagName string

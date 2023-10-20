package logus

/*
Structured logging with slog. To provide more rich logging info.
Slightly modified for comfort

!NOTE: first msg is always text.
They keys and values are going in pairs//

Optionally as single value can be added slogGroup
*/

import (
	"autogit/settings/envs"
	"autogit/settings/types"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

var (
	slogger *slog.Logger
)

func Debug(msg string, opts ...slogParam) {
	args := append([]any{}, newSlogGroup(opts...))
	if LOG_SHOW_FILE_LOCATIONS {
		args = append(args, logGroupFiles())
	}
	slogger.Debug(msg, args...)
}

func Info(msg string, opts ...slogParam) {
	args := append([]any{}, newSlogGroup(opts...))
	if LOG_SHOW_FILE_LOCATIONS {
		args = append(args, logGroupFiles())
	}
	slogger.Info(msg, args...)
}

func Warn(msg string, opts ...slogParam) {
	args := append([]any{}, newSlogGroup(opts...))
	if LOG_SHOW_FILE_LOCATIONS {
		args = append(args, logGroupFiles())
	}
	slogger.Warn(msg, args...)
}

func Error(msg string, opts ...slogParam) {
	args := append([]any{}, newSlogGroup(opts...))
	if LOG_SHOW_FILE_LOCATIONS {
		args = append(args, logGroupFiles())
	}
	slogger.Error(msg, args...)
}

func Fatal(msg string, opts ...slogParam) {
	args := append([]any{}, newSlogGroup(opts...))
	if LOG_SHOW_FILE_LOCATIONS {
		args = append(args, logGroupFiles())
	}
	slogger.Error(msg, args...)
	panic(msg)
}

func CheckError(err error, msg string, opts ...slogParam) {
	if err == nil {
		return
	}
	args := append([]any{}, newSlogGroup(opts...))
	args = append(args, "error")
	args = append(args, fmt.Sprintf("%v", err))
	slogger.Error(msg, args...)
	os.Exit(1)
}

func CheckFatal(err error, msg string, opts ...slogParam) {
	if err == nil {
		return
	}
	args := append([]any{}, newSlogGroup(opts...))
	args = append(args, "error")
	args = append(args, fmt.Sprintf("%v", err))
	slogger.Error(msg, args...)
	panic(msg)
}

func Debugf(msg string, varname string, value any, opts ...slogParam) {
	args := append([]any{}, newSlogGroup(opts...))
	if LOG_SHOW_FILE_LOCATIONS {
		args = append(args, logGroupFiles())
	}
	args = append(args, varname)
	args = append(args, fmt.Sprintf("%v", value))
	slogger.Debug(msg, args...)
}

func Infof(msg string, varname string, value any, opts ...slogParam) {
	args := append([]any{}, newSlogGroup(opts...))
	if LOG_SHOW_FILE_LOCATIONS {
		args = append(args, logGroupFiles())
	}
	args = append(args, varname)
	args = append(args, fmt.Sprintf("%v", value))
	slogger.Info(msg, args...)
}

var LOG_SHOW_FILE_LOCATIONS bool

func init() {
	LOG_SHOW_FILE_LOCATIONS = envs.LogShowFileLocations

	slogger = NewLogger(envs.LogLevel)
}

const (
	LEVEL_DEBUG types.LogLevel = "DEBUG"
	LEVEL_INFO  types.LogLevel = "INFO"
	LEVEL_WARN  types.LogLevel = "WARN"
	LEVEL_ERROR types.LogLevel = "ERROR"
)

func NewLogger(log_level_str types.LogLevel) *slog.Logger {
	var programLevel = new(slog.LevelVar) // Info by default

	switch log_level_str {
	case LEVEL_DEBUG:
		programLevel.Set(slog.LevelDebug)
	case LEVEL_INFO:
		programLevel.Set(slog.LevelInfo)
	case LEVEL_WARN:
		programLevel.Set(slog.LevelWarn)
	case LEVEL_ERROR:
		programLevel.Set(slog.LevelError)
	}

	if envs.LogTurnJSONLogging {
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
}

func GetCallingFile(level int) string {
	GetTwiceParentFunctionLocation := level
	_, filename, _, _ := runtime.Caller(GetTwiceParentFunctionLocation)
	filename = filepath.Base(filename)
	return fmt.Sprintf("f:%s ", filename)
}

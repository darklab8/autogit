package logus

/*
Structured logging with slog. To provide more rich logging info.
Slightly modified for comfort

!NOTE: first msg is always text.
They keys and values are going in pairs//

Optionally as single value can be added slogGroup
*/

import (
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

func CheckFatal(err error, msg string, opts ...slogParam) {
	if err == nil {
		return
	}
	args := append([]any{}, newSlogGroup(opts...))
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

var LOG_SHOW_FILE_LOCATIONS = false

func init() {
	show_files := os.Getenv("AUTOGIT_LOG_SHOW_FILE_LOCATIONS")
	LOG_SHOW_FILE_LOCATIONS = (show_files == "true")

	log_level_str, is_log_level_set := os.LookupEnv("AUTOGIT_LOG_LEVEL")
	if !is_log_level_set {
		log_level_str = "INFO"
	}
	slogger = NewLogger(LogLevel(log_level_str))
}

type LogLevel string

const (
	LEVEL_DEBUG LogLevel = "DEBUG"
	LEVEL_INFO  LogLevel = "INFO"
	LEVEL_WARN  LogLevel = "WARN"
	LEVEL_ERROR LogLevel = "ERROR"
)

func NewLogger(log_level_str LogLevel) *slog.Logger {
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

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
}

func GetCallingFile(level int) string {
	GetTwiceParentFunctionLocation := level
	_, filename, _, _ := runtime.Caller(GetTwiceParentFunctionLocation)
	filename = filepath.Base(filename)
	return fmt.Sprintf("f:%s ", filename)
}

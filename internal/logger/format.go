package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func Debugf(format string, args ...any) {
	slog.Debug(fmt.Sprintf(format, args...))
}

func Infof(format string, args ...any) {
	slog.Info(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...any) {
	slog.Warn(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...any) {
	slog.Error(fmt.Sprintf(format, args...))
}

func Fatal(msg string) {
	slog.Error(msg)
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	Errorf(format, args...)
	os.Exit(1)
}

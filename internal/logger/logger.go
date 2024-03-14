package logger

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"
)

type LogHandler struct {
	slog.Handler

	Attrs  []slog.Attr
	Groups []string
	Logger *log.Logger
}

func NewLogger(w io.Writer, opts *slog.HandlerOptions) *slog.Logger {
	handler := &LogHandler{
		Handler: slog.NewTextHandler(w, opts),
		Logger:  log.New(w, "", 0),
	}
	return slog.New(handler)
}

func ParseLevel(lvl string) (slog.Level, error) {
	switch strings.ToLower(lvl) {
	case "panic", "fatal", "error":
		return slog.LevelError, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "info":
		return slog.LevelInfo, nil
	case "debug", "trace":
		return slog.LevelDebug, nil
	}

	var l slog.Level
	return l, fmt.Errorf("%q is not a valid slog level", lvl)
}

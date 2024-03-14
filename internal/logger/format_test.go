package logger_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/Thibaut-gauvin/kie/internal/logger"
)

func Test_Format(t *testing.T) {
	t.Parallel()

	slog.SetDefault(logger.NewLogger(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Debugf("should not be displayed")

	slog.SetDefault(logger.NewLogger(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Debugf("should be displayed")

	logger.Infof("this is info")
	logger.Warnf("this is a warning")
	logger.Errorf("this is an error")
}

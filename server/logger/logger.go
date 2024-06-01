package logger

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() { //nolint:gochecknoinits
	logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

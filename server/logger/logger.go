package logger

import (
	"go.uber.org/zap"
)

var (
	Sugar *zap.SugaredLogger
)

func init() { //nolint:gochecknoinits
	log, _ := zap.NewProduction()
	Sugar = log.Sugar()
}

func Info(message string, fields ...any) {
	Sugar.Infof(message, fields)
}

func Error(message string, fields ...any) {
	Sugar.Errorf(message, fields)
}

func Fatal(message string, fields ...any) {
	Sugar.Fatalf(message, fields)
}

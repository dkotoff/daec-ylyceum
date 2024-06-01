package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.Out = os.Stdout

	logger.SetLevel(logrus.DebugLevel)
}

func Info(msg string, args ...any) {
	logger.WithFields(logrus.Fields{"Service": "Agent"}).Infof(msg, args...)
}

func Debug(msg string, args ...any) {
	logger.WithFields(logrus.Fields{"Service": "Agent"}).Debugf(msg, args...)
}

func Error(msg string, args ...any) {
	logger.WithFields(logrus.Fields{"Service": "Agent"}).Errorf(msg, args...)
}

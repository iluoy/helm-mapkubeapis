package log

import (
	"github.com/sirupsen/logrus"
)

//const logTimestampFormat = "2023-06-29 15:04:05"

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
		//		TimestampFormat: logTimestampFormat,
	})

	return logger
}

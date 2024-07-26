package main

import (
	"github.com/VikaPaz/message_server/internal/app"
	"github.com/sirupsen/logrus"
)

// @title Messages server API
// @description This is messages_server server.
// @host localhost:8800
func main() {

	logger := NewLogger(logrus.DebugLevel, &logrus.TextFormatter{
		FullTimestamp: true,
	})

	err := app.Run(logger)
	if err != nil {
		logger.Fatalln(err)
	}
}

func NewLogger(level logrus.Level, formatter logrus.Formatter) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(formatter)
	return logger
}

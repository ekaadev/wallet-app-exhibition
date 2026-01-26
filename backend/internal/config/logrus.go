package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NewLogger creates and returns a new Logrus logger instance configured for the application.
// Logger for logging application.
func NewLogger(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(viper.GetInt("log.level")))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}

package logging

import (
	"github.com/sirupsen/logrus"
)

const separator = "\n\n\n"

// Logger is an instance of the logrus type logger
var (
	Logger = logrus.StandardLogger()
)

func init() {
	Logger.SetLevel(logrus.TraceLevel)
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.Trace(separator)
}

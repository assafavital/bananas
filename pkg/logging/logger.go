package logging

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Logger is an instance of the logrus type logger
var (
	Logger = logrus.StandardLogger()
)

func NewFileLogger(logFileLevel logrus.Level, logFile string) *logrus.Logger {
	logger := logrus.New()
	// Send all logs to nowhere by default
	logger.SetOutput(io.Discard)
	// level is filtered by the writer hooks anyway...
	logger.SetLevel(logrus.TraceLevel)
	rotate := Rotate(logFile)
	logger.AddHook(&WriterHook{
		Writer: rotate,
		Level:  logFileLevel,
		Formatter: &withCMDFormatter{
			&stackTraceFormatter{
				&logrus.JSONFormatter{},
			},
		},
	})
	return logger
}

func Setup(logFileLevel logrus.Level, logFile string) {
	Logger = NewFileLogger(logFileLevel, logFile)
}

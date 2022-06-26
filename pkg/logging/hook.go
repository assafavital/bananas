package logging

import (
	"io"

	"github.com/sirupsen/logrus"
)

// code "borrowed" from https://github.com/sirupsen/logrus/issues/678
// WriterHook is a hook that writes logs of specified LogLevels to specified Writer
type WriterHook struct {
	Writer    io.Writer
	Level     logrus.Level
	Formatter logrus.Formatter
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterHook) Levels() []logrus.Level {
	levels := []logrus.Level{}
	for _, level := range logrus.AllLevels {
		if level <= hook.Level {
			levels = append(levels, level)
		}
	}
	return levels
}

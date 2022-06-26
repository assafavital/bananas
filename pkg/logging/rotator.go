package logging

import (
	"io"

	"gopkg.in/natefinch/lumberjack.v2"
)

func Rotate(path string) io.Writer {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    10,
		MaxAge:     30,
		MaxBackups: 5,
		LocalTime:  true,
		Compress:   true,
	}
}

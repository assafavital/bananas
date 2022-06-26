package config

import (
	"path/filepath"
)

const (
	DefaultLogFileName = "raftt.log"
	remoteLogDirPath   = "/var/log/raftt"
)

func LogFilePath() string {
	return filepath.Join(remoteLogDirPath, DefaultLogFileName)
}

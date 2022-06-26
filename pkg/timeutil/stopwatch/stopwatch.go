package stopwatch

import (
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	l "raftt.io/bananas/pkg/logging"
)

// Stopwatch logs the start and stop times with the specified message.
// use with either `defer timeutil.Stopwatch(...)()`
// or explicitly: `stop := timeutil.Stopwatch(...); ...; stop()`
// Note: reasoning for ellipsis in fieldsArr is to make it optional.
func Stopwatch(msg string, fieldsArr ...logrus.Fields) (stop func()) {
	instanceID := uuid.Must(uuid.NewRandom())
	logger := l.Logger.WithFields(logrus.Fields{
		"stopwatchKey": instanceID,
		"watching":     msg,
	})
	for _, fields := range fieldsArr {
		logger = logger.WithFields(fields)
	}

	logger.Debugf("Starting stopwatch: %s", msg)

	start := time.Now()

	return func() {
		logger.WithField(l.TimingField, time.Since(start)).Debugf("Finished stopwatch: %s", msg)
	}
}

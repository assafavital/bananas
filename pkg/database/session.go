package database

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Session struct {
	*gorm.DB
	URI string
}

// Same as gorm default logger, but ignoring record not found.
var defaultLogger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
	SlowThreshold:             200 * time.Millisecond,
	LogLevel:                  logger.Warn,
	IgnoreRecordNotFoundError: true,
	Colorful:                  true,
})

func MakeSession(uri string) (*Session, error) {
	return MakeSessionWithLogger(uri, defaultLogger)
}

func MakeSessionWithLogger(uri string, newLogger logger.Interface) (*Session, error) {
	session := Session{URI: uri}
	if err := session.initialize(newLogger); err != nil {
		return nil, err
	}
	return &session, nil
}

func (session *Session) initialize(newLogger logger.Interface) (err error) {
	session.DB, err = gorm.Open(postgres.Open(session.URI), &gorm.Config{Logger: newLogger})
	return errors.WithStack(err)
}

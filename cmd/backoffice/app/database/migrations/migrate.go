package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"raftt.io/bananas/cmd/backoffice/app/database/session"
)

// Migrate updates DB to current schema.
func Migrate() error {
	session, err := session.CurrentSession()
	if err != nil {
		return err
	}
	return MigrateConnection(session.DB)
}

func MigrateConnection(connection *gorm.DB) error {
	migrator := gormigrate.New(
		connection, gormigrate.DefaultOptions,
		append(
			[]*gormigrate.Migration{},
		),
	)
	return errors.WithStack(migrator.Migrate())
}

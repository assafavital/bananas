package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	l "raftt.io/bananas/pkg/logging"
)

func Migrate(db *gorm.DB, migrations []*gormigrate.Migration) error {
	l.Logger.Debug("🛂 Running migrations...")
	migrator := gormigrate.New(db, gormigrate.DefaultOptions, migrations)
	return errors.WithStack(migrator.Migrate())
}

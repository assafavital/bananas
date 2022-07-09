package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	l "raftt.io/bananas/pkg/logging"
)

func AddPurpleHair() *gormigrate.Migration {
	type PrimeMinister struct {
		HasPurpleHair bool
	}
	return &gormigrate.Migration{
		ID: "Add purple hair column to prime ministers table",
		Migrate: func(db *gorm.DB) error {
			l.Logger.Info("😈 Running purple hair migration")
			return errors.WithStack(db.AutoMigrate(&PrimeMinister{}))
		},
	}
}

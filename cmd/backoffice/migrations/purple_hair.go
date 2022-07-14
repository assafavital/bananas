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
			if err := db.AutoMigrate(&PrimeMinister{}); err != nil {
				return errors.WithStack(err)
			}

			// TODO: PMs named "bibi" should ALWAYS have purple hair!
			res := db.Table("prime_ministers").Where("name = ?", "bibi").Update("has_purple_hair", true)
			return res.Error
		},
	}
}

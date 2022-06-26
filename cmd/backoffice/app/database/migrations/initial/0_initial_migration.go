package initial

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	l "raftt.io/bananas/pkg/logging"
)

type User struct {
	gorm.Model
	Name string
}

func InitialMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0 Initial migration",
		Migrate: func(tx *gorm.DB) error {
			l.Logger.Info("Running first migration, creating empty tables")
			return errors.WithStack(tx.AutoMigrate(&User{}))
		},
	}
}

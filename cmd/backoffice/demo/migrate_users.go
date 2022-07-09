package demo

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"raftt.io/bananas/cmd/backoffice/migrations"
)

func MigrateUsers(db *gorm.DB) error {
	return migrations.Migrate(db, []*gormigrate.Migration{
		migrations.EnhanceUsersAndTokens(),
	})
}

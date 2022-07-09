package demo

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"raftt.io/bananas/cmd/backoffice/migrations"
	"raftt.io/bananas/cmd/backoffice/migrations/tokens"
)

func BetterMigrateUsers(db *gorm.DB) error {
	return migrations.Migrate(db, []*gormigrate.Migration{
		tokens.CreateTokensTable(),
		tokens.FetchAccountIDs(),
		tokens.DropTokenColumns(),
	})
}

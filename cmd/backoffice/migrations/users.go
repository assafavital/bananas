package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"raftt.io/bananas/cmd/backoffice/types"
	l "raftt.io/bananas/pkg/logging"
)

func CreateUsersTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "Create users table",
		Migrate: func(db *gorm.DB) error {
			l.Logger.Debug("✅ Creating 'users' table...")
			return db.AutoMigrate(&types.User{})
		},
	}
}

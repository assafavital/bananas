package tokens

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	l "raftt.io/bananas/pkg/logging"
)

func DropTokenColumns() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "Dropping token columns from users table",
		Migrate: func(db *gorm.DB) error {
			l.Logger.Debug("❌ Removing tokenless users and dropping token columns from 'users' table...")
			if err := deleteTokenlessUsers(db); err != nil {
				return err
			}
			return dropTokenColumnsFromUsers(db)
		},
	}
}

func deleteTokenlessUsers(db *gorm.DB) error {
	var usersWithTokens []uint
	if err := db.Table("user_tokens").Distinct("user_id").Find(&usersWithTokens).Error; err != nil {
		return err
	}
	return errors.WithStack(db.Delete(&User{}, usersWithTokens).Error)
}

func dropTokenColumnsFromUsers(db *gorm.DB) error {
	return db.Transaction(func(db *gorm.DB) error {
		if err := db.Migrator().DropColumn(&User{}, "github_token"); err != nil {
			return errors.WithStack(err)
		}
		if err := db.Migrator().DropColumn(&User{}, "gitlab_refresh_token"); err != nil {
			return errors.WithStack(err)
		}
		return errors.WithStack(db.Migrator().DropColumn(&User{}, "bibibucket_refresh_token"))
	})
}

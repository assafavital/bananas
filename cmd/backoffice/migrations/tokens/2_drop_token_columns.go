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
			return errors.WithStack(
				db.Exec(
					"ALTER TABLE users " +
						"DROP COLUMN github_token, " +
						"DROP COLUMN gitlab_refresh_token, " +
						"DROP COLUMN bibibucket_refresh_token",
				).Error,
			)
		},
	}
}

func deleteTokenlessUsers(db *gorm.DB) error {
	return errors.WithStack(db.Exec(
		"DELETE FROM users WHERE id IN (" +
			"SELECT users.id FROM USERS	LEFT OUTER JOIN user_tokens ON users.id = user_tokens.user_id " +
			"WHERE user_tokens.user_id IS NULL)").Error)
}

package demo

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"raftt.io/bananas/cmd/backoffice/migrations"
	"raftt.io/bananas/cmd/backoffice/types"
	l "raftt.io/bananas/pkg/logging"
)

func CreateUsers(db *gorm.DB) error {
	// Create users table
	if err := createUsersTable(db); err != nil {
		return err
	}

	// Create some users
	if err := createUsers(db); err != nil {
		return err
	}
	l.Logger.Debug("👥 Created users successfully!")
	return nil
}

func createUsersTable(db *gorm.DB) error {
	return migrations.Migrate(db,
		[]*gormigrate.Migration{
			migrations.CreateUsersTable(),
		},
	)
}

func createUsers(db *gorm.DB) error {
	users := []types.User{
		{
			Name:        "Marcus",
			GithubToken: "ghp_oi2h34t0",
		},
		{
			Name: "Papi",
		},
		{
			Name:                   "Assaf",
			GithubToken:            "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",
			BibibucketRefreshToken: "bbb_12345",
		},
	}
	return errors.Wrap(
		db.Create(&users).Error,
		"failed to create users",
	)
}

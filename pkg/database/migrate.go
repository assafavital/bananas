package database

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"raftt.io/bananas/pkg/timeutil"
)

const (
	databaseMigrationTimeout  = time.Second * 30
	databaseMigrationInterval = time.Second * 10
)

func Migrate(ctx context.Context, migrateFunc func() error) error {
	ctx, cancel := context.WithTimeout(ctx, databaseMigrationTimeout)
	defer cancel()
	return errors.Wrap(
		timeutil.Retry(ctx, databaseMigrationInterval, migrateFunc), "timeout waiting for migration")
}

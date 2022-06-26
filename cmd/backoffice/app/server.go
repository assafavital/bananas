package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"raftt.io/bananas/cmd/backoffice/app/database/migrations"
	"raftt.io/bananas/pkg/database"
)

const serverPort = 8081

var logger = logrus.StandardLogger()

func RunServer(ctx context.Context) error {
	logger.Info("Migrating the database...")
	if err := database.Migrate(ctx, migrations.Migrate); err != nil {
		return err
	}
	mux := http.NewServeMux()
	logger.Infof("Starting backoffice server at port %d", serverPort)
	return errors.WithStack(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), mux))
}

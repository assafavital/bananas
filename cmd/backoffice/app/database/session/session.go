package session

import (
	"os"

	"raftt.io/bananas/pkg/database"
)

func CurrentSession() (*database.Session, error) {
	return database.CurrentSession(os.Getenv("DB_URL"))
}

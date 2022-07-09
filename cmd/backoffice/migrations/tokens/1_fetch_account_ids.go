package tokens

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"raftt.io/bananas/cmd/backoffice/types"
	"raftt.io/bananas/pkg/bibibucket"
	"raftt.io/bananas/pkg/github"
	"raftt.io/bananas/pkg/gitlab"
	l "raftt.io/bananas/pkg/logging"
)

func FetchAccountIDs() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "Fetching account IDs",
		Migrate: func(db *gorm.DB) error {
			l.Logger.Debug("🧮 Fetching account IDs...")
			var tokens []types.Token
			if err := db.Model(&types.Token{}).Scan(&tokens).Error; err != nil {
				return err
			}
			for _, token := range tokens {
				if err := updateToken(db, token); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func updateToken(db *gorm.DB, token types.Token) error {
	var (
		accountID string
		err       error
	)

	switch token.Provider {
	case types.ProviderGithub:
		accountID, err = github.GetAccountID(token.Value)
	case types.ProviderGitlab:
		accountID, err = gitlab.GetAccountID(token.Value)
	case types.ProviderBibibucket:
		accountID, err = bibibucket.GetAccountID(token.Value)
	}

	if err != nil {
		l.Logger.WithError(err).
			WithField("token", token).
			Warn("failed to fetch account ID, removing invalid token")
		return errors.WithStack(db.Delete(&token).Error)
	}
	token.AccountID = accountID
	return errors.WithStack(db.Save(&token).Error)
}

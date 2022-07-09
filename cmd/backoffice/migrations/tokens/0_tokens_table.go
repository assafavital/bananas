package tokens

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"raftt.io/bananas/cmd/backoffice/types"
	l "raftt.io/bananas/pkg/logging"
)

type User struct {
	gorm.Model
	Name   string
	Tokens []types.Token `gorm:"many2many:user_tokens"`
}

func CreateTokensTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "Create and populate tokens table",
		Migrate: func(db *gorm.DB) error {
			l.Logger.Debug("🪙 Creating tokens from 'users' table...")
			if err := db.AutoMigrate(&User{}, &types.Token{}); err != nil {
				return errors.WithStack(err)
			}

			var users []types.User
			if err := db.Model(&types.User{}).Scan(&users).Error; err != nil {
				return err
			}
			for _, user := range users {
				if err := createUserTokens(db, user); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func createUserTokens(db *gorm.DB, user types.User) error {
	newUser := User{
		Model:  gorm.Model{ID: user.ID},
		Name:   user.Name,
		Tokens: make([]types.Token, 0),
	}

	if user.GithubToken != "" {
		newUser.Tokens = append(newUser.Tokens, getOrCreateToken(db, user.GithubToken, types.ProviderGithub))
	}

	if user.GitlabRefreshToken != "" {
		newUser.Tokens = append(newUser.Tokens, getOrCreateToken(db, user.GitlabRefreshToken, types.ProviderGitlab))
	}

	if user.BibibucketRefreshToken != "" {
		newUser.Tokens = append(newUser.Tokens, getOrCreateToken(db, user.BibibucketRefreshToken, types.ProviderBibibucket))
	}

	return errors.WithStack(db.Save(&newUser).Error)
}

func getOrCreateToken(db *gorm.DB, value string, provider types.TokenProvider) types.Token {
	if token, err := lookupTokenByValue(db, value); err == nil {
		return token
	}
	return types.Token{
		Provider:  provider,
		Value:     value,
		AccountID: types.EmptyAccountID,
	}
}

func lookupTokenByValue(tx *gorm.DB, value string) (types.Token, error) {
	var token types.Token
	err := tx.Model(&types.Token{}).Where("value = ?", value).First(&token).Error
	return token, errors.WithStack(err)
}

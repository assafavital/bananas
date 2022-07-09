package migrations

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

type User struct {
	gorm.Model
	Name   string
	Tokens []types.Token `gorm:"many2many:user_tokens"`
}

func EnhanceUsersAndTokens() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "Enhance users and tokens",
		Migrate: func(db *gorm.DB) error {
			l.Logger.Debug("✨ Enhancing users and tokens...")

			// fetch all users
			var users []types.User
			if err := db.Find(&users).Error; err != nil {
				return err
			}

			// remove tokens from users table
			if err := dropTokenColumnsFromUsers(db); err != nil {
				return err
			}

			// enhance all users
			for _, user := range users {
				if err := migrateUser(db, user); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func dropTokenColumnsFromUsers(db *gorm.DB) error {
	return errors.WithStack(
		db.Exec(
			"ALTER TABLE users " +
				"DROP COLUMN github_token, " +
				"DROP COLUMN gitlab_refresh_token, " +
				"DROP COLUMN bibibucket_refresh_token",
		).Error,
	)
}

func migrateUser(db *gorm.DB, user types.User) (err error) {
	newUser := User{
		Model:  gorm.Model{ID: user.ID},
		Name:   user.Name,
		Tokens: make([]types.Token, 3),
	}

	if newUser.Tokens[0], err = createToken(user.GithubToken, types.ProviderGithub); err != nil {
		return
	}
	if newUser.Tokens[1], err = createToken(user.GitlabRefreshToken, types.ProviderGitlab); err != nil {
		return
	}
	if newUser.Tokens[2], err = createToken(user.BibibucketRefreshToken, types.ProviderBibibucket); err != nil {
		return
	}

	return errors.WithStack(db.Save(&newUser).Error)
}

func createToken(token string, provider types.TokenProvider) (types.Token, error) {
	switch provider {
	case types.ProviderGithub:
		return createGithubToken(token)
	case types.ProviderGitlab:
		return createGitlabToken(token)
	case types.ProviderBibibucket:
		return createBibibucketToken(token)
	default:
		return types.Token{}, errors.New("no such token provider")
	}
}

func createGithubToken(githubToken string) (types.Token, error) {
	accountID, err := github.GetAccountID(githubToken)
	if err != nil {
		return types.Token{}, err
	}
	return types.Token{
		Value:     githubToken,
		Provider:  types.ProviderGithub,
		AccountID: accountID,
	}, nil
}

func createGitlabToken(gitlabRefreshToken string) (types.Token, error) {
	accountID, err := gitlab.GetAccountID(gitlabRefreshToken)
	if err != nil {
		return types.Token{}, err
	}
	return types.Token{
		Value:     gitlabRefreshToken,
		Provider:  types.ProviderGitlab,
		AccountID: accountID,
	}, nil
}

func createBibibucketToken(bibibucketRefreshToken string) (types.Token, error) {
	accountID, err := bibibucket.GetAccountID(bibibucketRefreshToken)
	if err != nil {
		return types.Token{}, err
	}
	return types.Token{
		Value:     bibibucketRefreshToken,
		Provider:  types.ProviderBibibucket,
		AccountID: accountID,
	}, nil
}

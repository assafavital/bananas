package migrations

import (
	"os"
	"testing"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/stretchr/testify/require"
	"raftt.io/bananas/cmd/backoffice/database/session"
	"raftt.io/bananas/cmd/backoffice/migrations/tokens"
	"raftt.io/bananas/cmd/backoffice/types"
	cf "raftt.io/bananas/pkg/controlflow"
	"raftt.io/bananas/pkg/database"
)

func createEmptyDB() {
	rootDBURI := os.Getenv("DB_URL_ROOT_DATABASE")
	session, err := database.MakeSession(rootDBURI)
	cf.PanicIfErr(err)
	cf.PanicIfErr(session.Exec("DROP DATABASE IF EXISTS postgres;").Error)
	cf.PanicIfErr(session.Exec("CREATE DATABASE postgres;").Error)
}

func TestMigration(t *testing.T) {
	// Create empty DB
	createEmptyDB()
	session, err := session.CurrentSession()
	require.NoError(t, err)

	// Define migrations
	createUsersTable := CreateUsersTable()
	createTokensTable := tokens.CreateTokensTable()
	fetchAccountIDs := tokens.FetchAccountIDs()
	dropTokenColumns := tokens.DropTokenColumns()
	migrator := gormigrate.New(session.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		createUsersTable,
		createTokensTable,
		fetchAccountIDs,
		dropTokenColumns,
	})

	// Test first migration
	require.NoError(t, migrator.MigrateTo(createTokensTable.ID))

	// Create some users
	db := session.DB
	require.NoError(t, db.Save([]types.User{
		{
			Name:                   "user1",
			GithubToken:            "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",
			BibibucketRefreshToken: "this is not a real token",
		},
		{
			Name:        "user2",
			GithubToken: "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",
		},
		{
			Name: "user3",
		},
	}).Error)

	// Test second migration
	require.NoError(t, migrator.MigrateTo(createTokensTable.ID))
	expectedTokens := map[uint][]types.Token{
		1: {
			{Value: "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", Provider: types.ProviderGithub},
			{Value: "this is not a real token", Provider: types.ProviderBibibucket},
		},
		2: {{Value: "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", Provider: types.ProviderGithub}},
	}
	type user struct{ ID uint }
	var users []user
	require.NoError(t, db.Find(&users).Error)
	for _, user := range users {
		var tokens []types.Token
		require.NoError(t, session.Raw("SELECT value, provider, account_id FROM tokens "+
			"JOIN user_tokens ON tokens.id = user_tokens.token_id "+
			"WHERE user_id = ?", user.ID).Scan(&tokens).Error)
		cf.Assert(len(tokens) == len(expectedTokens[user.ID]), "wrong number of tokens for user")
		for i, token := range tokens {
			cf.Assert(token.AccountID == "***", "expecting account id to be ***")
			cf.Assert(tokens[i].Value == expectedTokens[user.ID][i].Value, "wrong token string")
			cf.Assert(tokens[i].Provider == expectedTokens[user.ID][i].Provider, "wrong token provider")
		}
	}

	// Test third migration
	require.NoError(t, migrator.MigrateTo(fetchAccountIDs.ID))
	expectedTokens = map[uint][]types.Token{
		1: {{Value: "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", Provider: types.ProviderGithub}},
		2: {{Value: "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", Provider: types.ProviderGithub}},
	}
	require.NoError(t, db.Find(&users).Error)
	for _, user := range users {
		var tokens []types.Token
		require.NoError(t, session.Raw("SELECT value, provider, account_id FROM tokens "+
			"JOIN user_tokens ON tokens.id = user_tokens.token_id "+
			"WHERE user_id = ?", user.ID).Scan(&tokens).Error)
		cf.Assert(len(tokens) == len(expectedTokens[user.ID]), "wrong number of tokens for user")
		for i, token := range tokens {
			cf.Assert(token.AccountID == "17474078", "expecting account id to be 17474078")
			cf.Assert(tokens[i].Value == expectedTokens[user.ID][i].Value, "wrong token string")
			cf.Assert(tokens[i].Provider == expectedTokens[user.ID][i].Provider, "wrong token provider")
		}
	}

	// Thid fourth migration
	require.NoError(t, migrator.MigrateTo(dropTokenColumns.ID))
	var foundUser user
	require.Error(t, db.Where("name = ?", "user3").First(&foundUser).Error)
	require.False(t, db.Migrator().HasColumn(&types.User{}, "github_token"))
}

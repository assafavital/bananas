package types

import "gorm.io/gorm"

type TokenProvider string

const (
	ProviderGithub     TokenProvider = "github"
	ProviderGitlab     TokenProvider = "gitlab"
	ProviderBibibucket TokenProvider = "bibibucket"
	EmptyAccountID                   = "***"
)

type Token struct {
	gorm.Model
	Value     string        `gorm:"check:value <> '';index:unique_token"`
	Provider  TokenProvider `gorm:"check:provider <> '';index:unique_token"`
	AccountID string        `gorm:"check:account_id <> ''"`
}

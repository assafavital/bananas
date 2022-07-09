package types

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name                   string
	GithubToken            string
	GitlabRefreshToken     string
	BibibucketRefreshToken string
}

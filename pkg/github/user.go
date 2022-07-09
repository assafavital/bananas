package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	token string
}

func (t tokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: t.token}, nil
}

func GetAccountID(token string) (string, error) {
	ctx := context.TODO()
	oauthClient := oauth2.NewClient(ctx, &tokenSource{
		token: token,
	})
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return "", errors.WithStack(err)
	}
	if user.ID == nil {
		return "", errors.New("got nil user ID")
	}
	return fmt.Sprint(user.GetID()), nil
}

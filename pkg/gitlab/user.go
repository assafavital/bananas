package gitlab

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func GetAccountID(accessToken string) (string, error) {
	client, err := gitlab.NewClient(accessToken)
	if err != nil {
		return "", errors.WithStack(err)
	}
	user, _, err := client.Users.CurrentUser(gitlab.WithToken(gitlab.OAuthToken, accessToken))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return fmt.Sprint(user.ID), nil
}

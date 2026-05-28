package auth0

import (
	"context"
	"fmt"

	"github.com/auth0/go-auth0/v2/authentication"
	"github.com/auth0/go-auth0/v2/authentication/database"
)

type Auth0InitAPI struct {
	authAPI             *authentication.Authentication
	auth0ConnectionType string
}

func NewAuth0InitAPI(domain, clientID, clientSecret, auth0Connection string) (*Auth0InitAPI, error) {

	authAPi, err := authentication.New(context.TODO(), domain, authentication.WithClientID(clientID), authentication.WithClientSecret(clientSecret))

	if err != nil {
		return nil, fmt.Errorf("error with auth0 API: %v", err)
	}

	return &Auth0InitAPI{
		authAPI:             authAPi,
		auth0ConnectionType: auth0Connection,
	}, nil
}

func (a *Auth0InitAPI) SignUp(email, password string) (*string, error) {
	userData := database.SignupRequest{
		Email:      email,
		Password:   password,
		Connection: a.auth0ConnectionType,
	}

	createdUser, err := a.authAPI.Database.Signup(context.Background(), userData)

	if err != nil {
		return nil, fmt.Errorf("error requesting user creation: %v", err)
	}
	return &createdUser.ID, nil
}

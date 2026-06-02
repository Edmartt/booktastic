package auth0

import (
	"context"
	"fmt"

	"github.com/auth0/go-auth0/v2/authentication"
	"github.com/auth0/go-auth0/v2/authentication/database"
	"github.com/auth0/go-auth0/v2/authentication/oauth"
)

type Auth0InitAPI struct {
	domain              string
	clientID            string
	clientSecret        string
	auth0Connection     string
	authAPI             *authentication.Authentication
	auth0ConnectionType string
	audience            string
}

func NewAuth0InitAPI(domain, clientID, clientSecret, auth0Connection, audience string) (*Auth0InitAPI, error) {

	authAPi, err := authentication.New(context.TODO(), domain, authentication.WithClientID(clientID), authentication.WithClientSecret(clientSecret))

	if err != nil {
		return nil, fmt.Errorf("error with auth0 API: %v", err)
	}

	return &Auth0InitAPI{
		domain:              domain,
		clientID:            clientID,
		clientSecret:        clientSecret,
		auth0Connection:     auth0Connection,
		authAPI:             authAPi,
		auth0ConnectionType: auth0Connection,
		audience:            audience,
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

func (a *Auth0InitAPI) Login(email, password string) (*string, error) {
	tokenSet, err := a.authAPI.OAuth.LoginWithPassword(context.Background(), oauth.LoginWithPasswordRequest{
		Username: email,
		Password: password,
		Audience: a.audience,
	}, oauth.IDTokenValidationOptions{})

	if err != nil {
		return nil, fmt.Errorf("error with login: %v", err)
	}

	return &tokenSet.AccessToken, nil
}

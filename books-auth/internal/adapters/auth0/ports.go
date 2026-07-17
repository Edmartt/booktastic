package auth0

import (
	"context"

	"github.com/auth0/go-auth0/v2/authentication"
	"github.com/auth0/go-auth0/v2/authentication/database"
	"github.com/auth0/go-auth0/v2/authentication/oauth"
)

type Auth0DatabaseAPI interface {
	Signup(ctx context.Context, params database.SignupRequest, opts ...authentication.RequestOption) (*database.SignupResponse, error)
}

type Auth0AuthAPI interface {
	LoginWithPassword(ctx context.Context, body oauth.LoginWithPasswordRequest, validationOptions oauth.IDTokenValidationOptions, opts ...authentication.RequestOption) (*oauth.TokenSet, error)
}

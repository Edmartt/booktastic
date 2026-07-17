package auth0

import (
	"context"
	"errors"
	"testing"

	"github.com/auth0/go-auth0/v2/authentication"
	"github.com/auth0/go-auth0/v2/authentication/database"
	"github.com/auth0/go-auth0/v2/authentication/oauth"
	"github.com/stretchr/testify/assert"
)

type mockAuthOAuth struct {
	isError bool
}

func (ma *mockAuthOAuth) LoginWithPassword(ctx context.Context, body oauth.LoginWithPasswordRequest, validationOptions oauth.IDTokenValidationOptions, opts ...authentication.RequestOption) (*oauth.TokenSet, error) {

	if ma.isError {
		return nil, errors.New("login failed")
	}
	tokenSetData := oauth.TokenSet{
		AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyNDI2MjIsImVtYWlsIjoiam9obi5kb2VAZXhhbXBsZS5jb20iLCJyb2xlIjoidXNlciJ9.q6sFm8v7qTqU4RJ9v3wFhL2nJGKxZ5R2m8D3V7Y4nXo",
	}

	return &tokenSetData, nil
}

type mockAuthDatabase struct {
	isError bool
}

func (mdb *mockAuthDatabase) Signup(ctx context.Context, params database.SignupRequest, opts ...authentication.RequestOption) (*database.SignupResponse, error) {
	if mdb.isError {
		return nil, errors.New("auth0 signup failed")
	}
	userData := database.SignupResponse{
		ID: "auth0|12345678",
	}

	return &userData, nil
}

func TestSignupOk(t *testing.T) {

	mockDB := &mockAuthDatabase{isError: false}
	mockOAuth := &mockAuthOAuth{}

	auth0Service := &Auth0InitAPI{
		authDatabase:        mockDB,
		authOAuth:           mockOAuth,
		auth0ConnectionType: "Username-Password-Authentication",
	}

	email := "mockemail@mail.com"
	password := "testpassword1*"

	userdata, err := auth0Service.SignUp(email, password)

	assert.NotNil(t, userdata)
	assert.Nil(t, err)

	assert.Equal(t, "auth0|12345678", *userdata)
}

func TestSignupFail(t *testing.T) {
	mockDB := &mockAuthDatabase{isError: true}
	mockOAuth := &mockAuthOAuth{}

	auth0Service := &Auth0InitAPI{
		authDatabase:        mockDB,
		authOAuth:           mockOAuth,
		auth0ConnectionType: "Username-Password-Authentication",
	}

	email := "mockemail@mail.com"
	password := "testpassword1*"

	userData, err := auth0Service.SignUp(email, password)

	assert.Nil(t, userData)
	assert.NotNil(t, err)
	assert.Error(t, err)

	assert.EqualError(t, err, "error requesting user creation: auth0 signup failed")
}

func TestLoginSuccess(t *testing.T) {
	mockDB := &mockAuthDatabase{isError: false}
	mockOAuth := &mockAuthOAuth{isError: false}

	auth0Service := &Auth0InitAPI{
		authDatabase:        mockDB,
		authOAuth:           mockOAuth,
		auth0ConnectionType: "Username-Password-Authentication",
	}

	email := "mockemail@mail.com"
	password := "testpassword1*"

	accessToken, err := auth0Service.Login(email, password)

	assert.NotNil(t, accessToken)
	assert.Nil(t, err)

	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyNDI2MjIsImVtYWlsIjoiam9obi5kb2VAZXhhbXBsZS5jb20iLCJyb2xlIjoidXNlciJ9.q6sFm8v7qTqU4RJ9v3wFhL2nJGKxZ5R2m8D3V7Y4nXo", *accessToken)
}

func TestLoginFail(t *testing.T) {
	mockDB := &mockAuthDatabase{isError: false}
	mockOAuth := &mockAuthOAuth{isError: true}

	auth0Service := &Auth0InitAPI{
		authDatabase:        mockDB,
		authOAuth:           mockOAuth,
		auth0ConnectionType: "Username-Password-Authentication",
	}

	email := "mockemail@mail.com"
	password := "testpassword1*"

	accessToken, err := auth0Service.Login(email, password)

	assert.NotNil(t, err)
	assert.Nil(t, accessToken)

	assert.EqualError(t, err, "error with login: login failed")
}

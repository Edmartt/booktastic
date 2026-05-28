package auth0

import (
	"fmt"
	"os"
)

type Auth0Config struct {
	Domain          string
	Audience        string
	ClientID        string
	ClientSecret    string
	Auth0Connection string
}

func LoadAuthConfig() (*Auth0Config, error) {
	domain := os.Getenv("AUTH0_DOMAIN")

	if domain == "" {
		return nil, fmt.Errorf("AUTH0_DOMAIN environment variable required")
	}

	audience := os.Getenv("AUTH0_AUDIENCE")

	if audience == "" {
		return nil, fmt.Errorf("AUTH0_AUDIENCE environment variable required")
	}

	clientID := os.Getenv("AUTH0_CLIENT_ID")

	if clientID == "" {
		return nil, fmt.Errorf("AUTH0_CLIENT_ID environment required")
	}
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	if clientSecret == "" {
		return nil, fmt.Errorf("AUTH0_CLIENT_SECRET environment required")
	}

	auth0ConnectionType := os.Getenv("AUTH0_CONNECTION")

	if auth0ConnectionType == "" {
		return nil, fmt.Errorf("AUTH0_CONNECTION environment required")
	}

	return &Auth0Config{
		Domain:          domain,
		Audience:        audience,
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		Auth0Connection: auth0ConnectionType,
	}, nil
}

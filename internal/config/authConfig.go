package config

import (
	"fmt"
	"os"
)

type Auth0Config struct {
	Domain   string
	Audience string
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

	return &Auth0Config{
		Domain:   domain,
		Audience: audience,
	}, nil
}

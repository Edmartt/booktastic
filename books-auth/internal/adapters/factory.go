package adapters

import (
	"fmt"

	"github.com/edmartt/booktastic-auth-service/internal/adapters/auth0"
	"github.com/edmartt/booktastic-auth-service/internal/core/ports"
)

type AuthProviderBundle struct {
	Provider  ports.AuthProvider
	Validator ports.TokenValidator
}

func NewAuthProvider(provider string) (*AuthProviderBundle, error) {
	switch provider {
	case "auth0":
		config, err := auth0.LoadAuthConfig()

		if err != nil {
			return nil, err
		}

		validator, err := auth0.NewAuth0TokenValidator(config.Domain, config.Audience)

		if err != nil {
			return nil, err
		}

		auth0Provider, err := auth0.NewAuth0InitAPI(config.Domain, config.ClientID, config.ClientSecret, config.Auth0Connection, config.Audience)

		if err != nil {
			return nil, err
		}
		return &AuthProviderBundle{
			Provider:  auth0Provider,
			Validator: validator,
		}, nil
	default:
		return nil, fmt.Errorf("unkown auth provider: %s", provider)
	}

}

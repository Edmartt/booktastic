package auth0

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/auth0/go-jwt-middleware/v3/jwks"
	"github.com/auth0/go-jwt-middleware/v3/validator"
)

func JWTvalidator(domain, audience string) (*validator.Validator, error) {
	issuerURL, err := url.Parse(os.Getenv("AUTH0_DOMAIN"))

	if err != nil {
		return nil, fmt.Errorf("failed to parse issuer URL")
	}
	provider, err := jwks.NewCachingProvider(

		jwks.WithIssuerURL(issuerURL),
		jwks.WithCacheTTL(5*time.Minute),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS provider: %w", err)
	}

	jwtValidator, err := validator.New(

		validator.WithKeyFunc(provider.KeyFunc),
		validator.WithAlgorithm(validator.RS256),
		validator.WithIssuer(issuerURL.String()),
		validator.WithAudience(audience),
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &CustomClaims{}
		}),
		validator.WithAllowedClockSkew(30*time.Second),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create validator: %w", err)
	}

	return jwtValidator, nil

}

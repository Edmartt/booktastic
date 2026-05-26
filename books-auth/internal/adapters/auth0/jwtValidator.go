package auth0

import (
	"context"
	"net/url"
	"time"

	"github.com/auth0/go-jwt-middleware/v3/jwks"
	"github.com/auth0/go-jwt-middleware/v3/validator"
	projectDomain "github.com/edmartt/booktastic-auth-service/internal/core/domain"
)

type Auth0TokenValidator struct {
	localValidator *validator.Validator
}

func NewAuth0TokenValidator(domain, audience string) (*Auth0TokenValidator, error) {
	issuerURL, _ := url.Parse(domain)
	provider, _ := jwks.NewCachingProvider(
		jwks.WithIssuerURL(issuerURL),
		jwks.WithCacheTTL(5*time.Minute),
	)

	jwtValidator, _ := validator.New(

		validator.WithKeyFunc(provider.KeyFunc),
		validator.WithAlgorithm(validator.RS256),
		validator.WithIssuer(issuerURL.String()),
		validator.WithAudience(audience),
		validator.WithCustomClaims(func() validator.CustomClaims {

			return &projectDomain.CustomClaims{}
		}),
		validator.WithAllowedClockSkew(30*time.Second),
	)

	return &Auth0TokenValidator{
		localValidator: jwtValidator,
	}, nil
}

func (v *Auth0TokenValidator) Validate(ctx context.Context, token string) (any, error) {

	return v.localValidator.ValidateToken(ctx, token)
}

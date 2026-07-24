package auth0

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCustomClaimsNoError(t *testing.T) {
	customClaims := CustomClaims{
		Scope: "read:books write:books",
	}

	claimsValidationResponse := customClaims.Validate(context.Background())

	assert.Nil(t, claimsValidationResponse, "unexpected error validating claims")
}
func TestValidateCustomClaimsScopeEmpty(t *testing.T) {
	customClaims := CustomClaims{
		Scope: "",
	}
	err := customClaims.Validate(context.Background())

	assert.Nil(t, err)
}

func TestValidateCustomClaimsInvalidWhitespace(t *testing.T) {
	customClaims := CustomClaims{Scope: " read:books"}
	err := customClaims.Validate(context.Background())
	assert.EqualError(t, err, "scope claim has invalid whitespace")
}
func TestValidateCustomClaimsDoubleWhitespace(t *testing.T) {
	customClaims := CustomClaims{Scope: "read:books  write:books"}
	err := customClaims.Validate(context.Background())
	assert.EqualError(t, err, "scope claim contains double spaces")
}

func TestHasScopeTrue(t *testing.T) {
	expectedScope := "read:books"

	cClaims := CustomClaims{
		Scope: "read:books write:books",
	}

	resultScope := cClaims.HasScope(expectedScope)

	assert.True(t, resultScope)
}

func TestHasScopeFalse(t *testing.T) {
	expectedScope := "read:books"

	cClaims := CustomClaims{
		Scope: "",
	}

	resultScope := cClaims.HasScope(expectedScope)

	assert.False(t, resultScope)
}

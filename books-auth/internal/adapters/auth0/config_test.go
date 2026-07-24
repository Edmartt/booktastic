package auth0

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAuthConfigEmptyDomain(t *testing.T) {

	configObject, err := LoadAuthConfig()

	assert.Error(t, err)

	assert.Nil(t, configObject)
}

func TestLoadAuthConfigEmptyAudience(t *testing.T) {

	os.Setenv("AUTH0_DOMAIN", "custom_domain")
	configObject, err := LoadAuthConfig()

	assert.Error(t, err)

	assert.Nil(t, configObject)
}

func TestLoadAuthConfigEmptyClientID(t *testing.T) {

	os.Setenv("AUTH0_AUDIENCE", "custom_audience")
	configObject, err := LoadAuthConfig()

	assert.Error(t, err)

	assert.Nil(t, configObject)
}

func TestLoadAuthConfigEmptyClientSecret(t *testing.T) {

	os.Setenv("AUTH0_CLIENT_ID", "custom_client_id")

	configObject, err := LoadAuthConfig()

	assert.Error(t, err)

	assert.Nil(t, configObject)
}

func TestLoadAuthConfigEmptyAuth0Connection(t *testing.T) {

	os.Setenv("AUTH0_CLIENT_SECRET", "custom_client_secret")

	configObject, err := LoadAuthConfig()

	assert.Error(t, err)

	assert.Nil(t, configObject)
}

func TestLoadAuthConfigOK(t *testing.T) {

	os.Setenv("AUTH0_AUDIENCE", "custom_audience")
	os.Setenv("AUTH0_DOMAIN", "custom_domain")
	os.Setenv("AUTH0_CLIENT_ID", "custom_client_secret")
	os.Setenv("AUTH0_CLIENT_SECRET", "custom_client_secret")
	os.Setenv("AUTH0_CONNECTION", "custom_auth0_connection")

	configObject, err := LoadAuthConfig()

	assert.NoError(t, err)

	assert.NotNil(t, configObject)

	assert.Equal(t, "custom_audience", os.Getenv("AUTH0_AUDIENCE"))
}

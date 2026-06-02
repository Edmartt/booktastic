package main

import (
	"os"

	"github.com/edmartt/booktastic-auth-service/internal/adapters/auth0"
	"github.com/edmartt/booktastic-auth-service/internal/adapters/http"
	"github.com/joho/godotenv"
)

func main() {

	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	}

	auth0Config, err := auth0.LoadAuthConfig()

	if err != nil {
		panic(err)
	}

	validator, err := auth0.NewAuth0TokenValidator(auth0Config.Domain, auth0Config.Audience)

	if err != nil {
		panic(err)
	}

	authProvider, err := auth0.NewAuth0InitAPI(auth0Config.Domain, auth0Config.ClientID, auth0Config.ClientSecret, auth0Config.Auth0Connection, auth0Config.Audience)

	if err != nil {
		panic(err)
	}

	handler := http.NewHandler(validator, authProvider)

	server := http.HTTPServer{
		Handler: *handler,
	}

	if err = server.RunServer(os.Getenv("HTTP_PORT")); err != nil {
		panic(err)
	}
}

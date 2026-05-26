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

	handler := http.NewHandler(validator)

	server := http.HTTPServer{
		Handler: *handler,
	}

	if err = server.RunServer(os.Getenv("HTTP_PORT")); err != nil {
		panic(err)
	}

}

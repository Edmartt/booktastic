package main

import (
	"os"

	"github.com/edmartt/booktastic-auth-service/internal/adapters"
	"github.com/edmartt/booktastic-auth-service/internal/adapters/http"
	errorHandling "github.com/edmartt/booktastic-shared/errors/http/adapters/ginhttp"
	"github.com/joho/godotenv"
)

// @title Booktastic Auth API
// @version 1.0
// @description Authentication service for Booktastic
// @host auth.localhost
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	}

	authProvider, err := adapters.NewAuthProvider(os.Getenv("AUTH_PROVIDER"))

	if err != nil {
		panic(err)
	}

	errorHandler := errorHandling.NewGinErrors()

	handler := http.NewHandler(authProvider.Validator, authProvider.Provider, *errorHandler)

	server := http.HTTPServer{
		Handler: *handler,
	}

	if err = server.RunServer(os.Getenv("HTTP_PORT")); err != nil {
		panic(err)
	}
}

package main

import (
	"os"

	"github.com/edmartt/bookstatic-book-service/internal/adapters/database"
	"github.com/edmartt/bookstatic-book-service/internal/adapters/http"
	"github.com/edmartt/bookstatic-book-service/internal/adapters/repository"
	errorHandling "github.com/edmartt/booktastic-shared/errors/http/adapters/ginhttp"
	"github.com/joho/godotenv"
)

// @title Booktastic Books API
// @version 1.0
// @description Books management service for Booktastic
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	}

	cfg, err := database.LoadPGConfig()

	if err != nil {
		panic(err)
	}
	dbConnectObject := database.NewPostgres(cfg)
	getConn := dbConnectObject.GetConnection()

	database.PingDB(getConn)
	db := repository.NewRepository(dbConnectObject)
	errorHandler := errorHandling.NewGinErrors()
	handlerObject := http.NewHandler(db, *errorHandler)
	server := http.HTTPServer{
		Handler: *handlerObject,
	}

	if err := server.RunServer(os.Getenv("HTTP_PORT")); err != nil {
		panic(err)
	}
}

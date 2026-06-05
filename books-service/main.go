package main

import (
	"os"

	"github.com/edmartt/bookstatic-book-service/internal/application"
	"github.com/edmartt/bookstatic-book-service/internal/books/data"
	"github.com/edmartt/bookstatic-book-service/internal/database"
	errorHandling "github.com/edmartt/booktastic-shared/errors/http/adapters/ginhttp"
	"github.com/joho/godotenv"
)

func main() {

	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	}

	dbConnectObject := &database.Postgres{}
	getConn := dbConnectObject.GetConnection()

	database.PingDB(getConn)
	db := data.NewRepository(dbConnectObject)
	errorHandler := errorHandling.NewGinErrors()
	handlerObject := application.NewHandler(db, *errorHandler)
	server := application.HTTPServer{
		Handler: *handlerObject,
	}

	if err := server.RunServer(os.Getenv("HTTP_PORT")); err != nil {
		panic(err)
	}
}

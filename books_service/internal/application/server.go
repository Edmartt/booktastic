package application

import (
	"fmt"
	"log"

	"github.com/auth0/go-jwt-middleware/v3/validator"
	"github.com/edmartt/bookstatic-book-service/internal/config"
	"github.com/edmartt/bookstatic-book-service/internal/utils"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	Handler HTTPHandler
}

func (h HTTPServer) setBookRoutes(router *gin.RouterGroup) {
	router.GET("/books/:id", h.Handler.ReadBook)
	router.POST("/books", h.Handler.CreateBook)
	router.PATCH("/books/:id", h.Handler.UpdateBook)
}

func (h HTTPServer) setRouter(validatorJWT *validator.Validator) *gin.Engine {
	router := gin.Default()

	router.Use(LimitRequest())

	router.Use(AuthMiddleware(validatorJWT))
	apiGroup := router.Group("/api/v1")
	h.setBookRoutes(apiGroup)

	return router
}

// RunServer starts http server
func (h HTTPServer) RunServer(port string) error {
	auth0ConfigObject, err := config.LoadAuthConfig()

	if err != nil {
		return fmt.Errorf("error loading AUTH0 Config: %w", err)
	}

	jwtValidator, err := utils.JWTvalidator(auth0ConfigObject.Domain, auth0ConfigObject.Audience)

	if err != nil {
		return fmt.Errorf("error creating validator: %w", err)
	}

	router := h.setRouter(jwtValidator)
	log.Fatal(router.Run(":" + port))

	return nil
}

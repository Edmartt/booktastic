package http

import (
	"log"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	Handler HTTPHandler
}

func (h HTTPServer) setAuthRoutes(router *gin.RouterGroup) {
	router.GET("/auth/verify", h.Handler.VerifyJWTToken)
	router.POST("/auth/signup", h.Handler.SignupUserHandler)
	router.POST("/auth/login", h.Handler.LoginUserHandler)
}

func (h HTTPServer) setRouter() *gin.Engine {
	router := gin.Default()

	apiGroup := router.Group("/api/v1")
	h.setAuthRoutes(apiGroup)

	return router
}

// RunServer starts http server
func (h HTTPServer) RunServer(port string) error {

	router := h.setRouter()
	log.Fatal(router.Run(":" + port))

	return nil
}

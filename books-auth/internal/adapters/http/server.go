package http

import (
	"log"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	Handler HTTPHandler
}

func (h HTTPServer) setBookRoutes(router *gin.RouterGroup) {
	router.GET("/auth/verify", h.Handler.VerifyJWTToken)
}

func (h HTTPServer) setRouter() *gin.Engine {
	router := gin.Default()

	apiGroup := router.Group("/api/v1")
	h.setBookRoutes(apiGroup)

	return router
}

// RunServer starts http server
func (h HTTPServer) RunServer(port string) error {

	router := h.setRouter()
	log.Fatal(router.Run(":" + port))

	return nil
}

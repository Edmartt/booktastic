package ginhttp

import "github.com/gin-gonic/gin"

type GinErrors struct {
}

func NewGinErrors() *GinErrors {
	return &GinErrors{}

}

func (g *GinErrors) WriteError(context *gin.Context, code int, message string) {
	context.AbortWithStatusJSON(
		code,
		gin.H{
			"message": message,
		},
	)
}

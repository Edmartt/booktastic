package ginhttp

import "github.com/gin-gonic/gin"

type GinErrors struct {
	ctx *gin.Context
}

func NewGinErrors(ctx *gin.Context) *GinErrors {
	return &GinErrors{
		ctx: ctx,
	}

}

func (g *GinErrors) WriteError(code int, message string) {
	g.ctx.AbortWithStatusJSON(
		code,
		gin.H{
			"message": message,
		},
	)
}

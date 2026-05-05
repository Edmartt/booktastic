package application

import (
	"github.com/edmartt/bookstatic-book-service/internal/utils"
	"github.com/gin-gonic/gin"
)

func LimitRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := utils.GetLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{"error": "too many requests"})
			return
		}

		c.Next()
	}
}

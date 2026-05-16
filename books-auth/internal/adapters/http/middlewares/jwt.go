package middlewares

import (
	"log/slog"
	"strings"

	"github.com/auth0/go-jwt-middleware/v3/validator"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtValidator *validator.Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			unauthorized(c, "Missing authorization header")
			return
		}

		const bearer = "Bearer "

		if !strings.HasPrefix(authHeader, bearer) {
			unauthorized(c, "Invalid authorization header")
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearer))

		claims, err := jwtValidator.ValidateToken(c.Request.Context(), token)
		if err != nil {
			slog.Error(
				"JWT validation failed",
				"error", err,
				"path", c.FullPath(),
			)

			unauthorized(c, "Failed to validate JWT")
			return
		}

		c.Set("jwt", claims)

		c.Next()
	}
}

func unauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(
		401,
		gin.H{
			"message": message,
		},
	)
}

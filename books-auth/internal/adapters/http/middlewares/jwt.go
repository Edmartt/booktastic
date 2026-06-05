package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/edmartt/booktastic-auth-service/internal/core/ports"
	selfErrors "github.com/edmartt/booktastic-shared/errors/http/adapters/ginhttp"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenValidator ports.TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		writer := selfErrors.NewGinErrors()

		if authHeader == "" {
			writer.WriteError(c, http.StatusUnauthorized, "missing authorization header")
			return
		}

		const bearer = "Bearer "

		if !strings.HasPrefix(authHeader, bearer) {
			writer.WriteError(c, http.StatusUnauthorized, "Invalid authorization header")
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearer))

		claims, err := tokenValidator.Validate(c.Request.Context(), token)

		if err != nil {
			slog.Error(
				"JWT validation failed",
				"error", err,
				"path", c.FullPath(),
			)

			writer.WriteError(c, http.StatusUnauthorized, "failed to validate JWT")
			return
		}

		c.Set("jwt", claims)

		c.Next()
	}
}

package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/auth0/go-jwt-middleware/v3/validator"
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	jwtValidator *validator.Validator
}

func NewHandler(jwtValidator *validator.Validator) *HTTPHandler {
	return &HTTPHandler{
		jwtValidator: jwtValidator,
	}
}

func (h HTTPHandler) VerifyJWTToken(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")

	if authHeader == "" {
		unauthorized(context, "Missing Authorization Header")
		return
	}

	const bearer = "Bearer"

	if !strings.HasPrefix(authHeader, bearer) {
		unauthorized(context, "Invalid authorization header")
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearer))

	_, err := h.jwtValidator.ValidateToken(context.Request.Context(), token)

	if err != nil {
		slog.Error("JWT validation failed", "error", err, "path", context.FullPath())
		unauthorized(context, "Failed to validate JWT")
		return
	}
}

func unauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": message,
	})
}

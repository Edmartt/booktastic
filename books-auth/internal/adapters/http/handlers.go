package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/edmartt/booktastic-auth-service/internal/adapters/http/dto"
	"github.com/edmartt/booktastic-auth-service/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	jwtValidator ports.TokenValidator
	authProvider ports.AuthProvider
}

func NewHandler(jwtValidator ports.TokenValidator, authProvider ports.AuthProvider) *HTTPHandler {
	return &HTTPHandler{
		jwtValidator: jwtValidator,
		authProvider: authProvider,
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

	_, err := h.jwtValidator.Validate(context.Request.Context(), token)

	if err != nil {
		slog.Error("JWT validation failed", "error", err, "path", context.FullPath())
		unauthorized(context, "Failed to validate JWT")
		return
	}
}

func (h HTTPHandler) SignupUserHandler(c *gin.Context) {
	var myDTO dto.SignUpDTO

	if err := c.BindJSON(&myDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error binding body"})
		return
	}

	if myDTO.PasswordConfirmation != myDTO.Password {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "password confirmation error"})
		return
	}

	if !myDTO.IsValidPassword() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "password must contain uppercase, lowercase, numbers and special characters"})
		return
	}

	authProviderResponse, err := h.authProvider.SignUp(myDTO.Email, myDTO.Password)

	slog.Error("unkown error", "error", err, "path", c.FullPath())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "user already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user created ID": *authProviderResponse})
}

func unauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": message,
	})
}

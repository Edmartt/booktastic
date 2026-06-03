package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/auth0/go-jwt-middleware/v3/validator"
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

	claims, err := h.jwtValidator.Validate(context.Request.Context(), token)

	if err != nil {
		slog.Error("JWT validation failed", "error", err, "path", context.FullPath())
		unauthorized(context, "Failed to validate JWT")
		return
	}

	validatedClaims := claims.(*validator.ValidatedClaims)

	context.Header("X-User-Id", validatedClaims.RegisteredClaims.Subject)
	context.Status(http.StatusOK)
}

func (h HTTPHandler) SignupUserHandler(context *gin.Context) {
	var myDTO dto.SignUpDTO

	if err := context.BindJSON(&myDTO); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error binding body"})
		return
	}

	if myDTO.PasswordConfirmation != myDTO.Password {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "password confirmation error"})
		return
	}

	if !myDTO.IsValidPassword() {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "password must contain uppercase, lowercase, numbers and special characters"})
		return
	}

	authProviderResponse, err := h.authProvider.SignUp(myDTO.Email, myDTO.Password)

	slog.Error("unkown error", "error", err, "path", context.FullPath())

	if err != nil {
		context.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "user already exists"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"user created ID": *authProviderResponse})
}

func unauthorized(context *gin.Context, message string) {
	context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": message,
	})
}

func (h *HTTPHandler) LoginUserHandler(context *gin.Context) {

	var loginDTO dto.LoginDTO

	if err := context.BindJSON(&loginDTO); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	auth0Token, err := h.authProvider.Login(loginDTO.Email, loginDTO.Password)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"access token ": auth0Token})

}

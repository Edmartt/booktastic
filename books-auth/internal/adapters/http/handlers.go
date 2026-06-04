package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/auth0/go-jwt-middleware/v3/validator"
	"github.com/edmartt/booktastic-auth-service/internal/adapters/http/dto"
	"github.com/edmartt/booktastic-auth-service/internal/core/ports"
	"github.com/gin-gonic/gin"

	selfErrors "github.com/edmartt/booktastic-shared/errors/http/adapters/ginhttp"
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
	writer := selfErrors.NewGinErrors(context)

	if authHeader == "" {
		writer.WriteError(401, "Missing Authorization Header")
		return
	}

	const bearer = "Bearer"

	if !strings.HasPrefix(authHeader, bearer) {
		writer.WriteError(401, "Invalid authorization header")
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearer))

	claims, err := h.jwtValidator.Validate(context.Request.Context(), token)

	if err != nil {
		slog.Error("JWT validation failed", "error", err, "path", context.FullPath())
		writer.WriteError(401, "failed to validate JWT")
		return
	}

	validatedClaims := claims.(*validator.ValidatedClaims)

	context.Header("X-User-Id", validatedClaims.RegisteredClaims.Subject)
	context.Status(http.StatusOK)
}

func (h HTTPHandler) SignupUserHandler(context *gin.Context) {
	var myDTO dto.SignUpDTO
	writer := selfErrors.NewGinErrors(context)

	if err := context.BindJSON(&myDTO); err != nil {
		writer.WriteError(http.StatusBadRequest, "error binding body")
		return
	}

	if myDTO.PasswordConfirmation != myDTO.Password {
		writer.WriteError(http.StatusBadRequest, "password confirmation error")
		return
	}

	if !myDTO.IsValidPassword() {
		writer.WriteError(http.StatusBadRequest, "password must contain uppercase, lowercase, numbers and special characters")
		return
	}

	authProviderResponse, err := h.authProvider.SignUp(myDTO.Email, myDTO.Password)

	slog.Error("unkown error", "error", err, "path", context.FullPath())

	if err != nil {
		writer.WriteError(http.StatusConflict, "user already exists")
		return
	}

	context.JSON(http.StatusCreated, gin.H{"user created ID": *authProviderResponse})
}

func (h *HTTPHandler) LoginUserHandler(context *gin.Context) {

	var loginDTO dto.LoginDTO
	writer := selfErrors.NewGinErrors(context)

	if err := context.BindJSON(&loginDTO); err != nil {
		writer.WriteError(http.StatusBadRequest, "error with user or password, check data")
		return
	}

	auth0Token, err := h.authProvider.Login(loginDTO.Email, loginDTO.Password)

	if err != nil {
		writer.WriteError(http.StatusUnauthorized, "invalid credentials")
		return
	}

	context.JSON(http.StatusOK, gin.H{"access token ": auth0Token})

}

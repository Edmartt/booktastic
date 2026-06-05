package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/auth0/go-jwt-middleware/v3/validator"
	"github.com/edmartt/booktastic-auth-service/internal/adapters/http/dto"
	"github.com/edmartt/booktastic-auth-service/internal/core/ports"
	"github.com/gin-gonic/gin"

	errorHandling "github.com/edmartt/booktastic-shared/errors/http/adapters/ginhttp"
)

type HTTPHandler struct {
	jwtValidator ports.TokenValidator
	authProvider ports.AuthProvider
	errorWriter  errorHandling.GinErrors
}

func NewHandler(jwtValidator ports.TokenValidator, authProvider ports.AuthProvider, errorWriter errorHandling.GinErrors) *HTTPHandler {
	return &HTTPHandler{
		jwtValidator: jwtValidator,
		authProvider: authProvider,
		errorWriter:  errorWriter,
	}
}

func (h HTTPHandler) VerifyJWTToken(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")

	if authHeader == "" {
		h.errorWriter.WriteError(context, http.StatusUnauthorized, "Missing Authorization Header")
		return
	}

	const bearer = "Bearer"

	if !strings.HasPrefix(authHeader, bearer) {
		h.errorWriter.WriteError(context, http.StatusUnauthorized, "Invalid authorization header")
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearer))

	claims, err := h.jwtValidator.Validate(context.Request.Context(), token)

	if err != nil {
		slog.Error("JWT validation failed", "error", err, "path", context.FullPath())
		h.errorWriter.WriteError(context, http.StatusUnauthorized, "failed to validate JWT")
		return
	}

	validatedClaims := claims.(*validator.ValidatedClaims)

	context.Header("X-User-Id", validatedClaims.RegisteredClaims.Subject)
	context.Status(http.StatusOK)
}

func (h HTTPHandler) SignupUserHandler(context *gin.Context) {
	var myDTO dto.SignUpDTO

	if err := context.BindJSON(&myDTO); err != nil {
		h.errorWriter.WriteError(context, http.StatusBadRequest, "error binding body")
		return
	}

	if myDTO.PasswordConfirmation != myDTO.Password {
		h.errorWriter.WriteError(context, http.StatusBadRequest, "password confirmation error")
		return
	}

	if !myDTO.IsValidPassword() {
		h.errorWriter.WriteError(context, http.StatusBadRequest, "password must contain uppercase, lowercase, numbers and special characters")
		return
	}

	authProviderResponse, err := h.authProvider.SignUp(myDTO.Email, myDTO.Password)

	slog.Error("unkown error", "error", err, "path", context.FullPath())

	if err != nil {
		h.errorWriter.WriteError(context, http.StatusConflict, "user already exists")
		return
	}

	context.JSON(http.StatusCreated, gin.H{"user created ID": *authProviderResponse})
}

func (h *HTTPHandler) LoginUserHandler(context *gin.Context) {

	var loginDTO dto.LoginDTO

	if err := context.BindJSON(&loginDTO); err != nil {
		h.errorWriter.WriteError(context, http.StatusBadRequest, "error with user or password, check data")
		return
	}

	auth0Token, err := h.authProvider.Login(loginDTO.Email, loginDTO.Password)

	if err != nil {
		h.errorWriter.WriteError(context, http.StatusUnauthorized, "invalid credentials")
		return
	}

	context.JSON(http.StatusOK, gin.H{"access token ": auth0Token})
}

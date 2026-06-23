package http

type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhdXRoMHwxMjM0NTY3ODkwIn0.signature"`
}

type UserCreatedResponse struct {
	UserCreatedID string `json:"user_created_id" example:"auth0|6a18b3add4d8019d2318f29b"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

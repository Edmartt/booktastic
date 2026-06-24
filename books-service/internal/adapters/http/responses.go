package http

type BookCreatedResponse struct {
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

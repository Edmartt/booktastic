package http

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/edmartt/bookstatic-book-service/internal/core/domain/models"
	"github.com/edmartt/bookstatic-book-service/internal/core/ports"
	errorHandling "github.com/edmartt/booktastic-shared/errors/http/adapters/ginhttp"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

// fakeBookRepository is the single in-memory implementation of
// ports.IDataAccessLayer shared by every handler test in this package. Each
// operation panics when its function is not configured so a misconfigured test
// fails loudly instead of simulating a handler response.
type fakeBookRepository struct {
	readFn   func(id, userID string) (*models.Books, error)
	createFn func(book models.Books) (string, error)
	updateFn func(book *models.Books) (*models.Books, error)
}

func (f fakeBookRepository) Create(book models.Books) (string, error) {
	if f.createFn == nil {
		panic("fakeBookRepository: createFn is not configured")
	}

	return f.createFn(book)
}

func (f fakeBookRepository) Read(id, userID string) (*models.Books, error) {
	if f.readFn == nil {
		panic("fakeBookRepository: readFn is not configured")
	}

	return f.readFn(id, userID)
}

func (f fakeBookRepository) Update(book *models.Books) (*models.Books, error) {
	if f.updateFn == nil {
		panic("fakeBookRepository: updateFn is not configured")
	}

	return f.updateFn(book)
}

func newTestHandler(repository ports.IDataAccessLayer) *HTTPHandler {
	errorWriter := errorHandling.NewGinErrors()

	return NewHandler(repository, *errorWriter)
}

func errorMessage(t *testing.T, recorder *httptest.ResponseRecorder) string {
	t.Helper()

	var response map[string]string
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))

	return response["message"]
}

func marshalPayload(t *testing.T, fields map[string]any) string {
	t.Helper()

	body, err := json.Marshal(fields)
	require.NoError(t, err)

	return string(body)
}

func stringPtr(value string) *string {
	return &value
}

func responseKeys(t *testing.T, recorder *httptest.ResponseRecorder) []string {
	t.Helper()

	var rawResponse map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &rawResponse))

	keys := make([]string, 0, len(rawResponse))
	for key := range rawResponse {
		keys = append(keys, key)
	}

	return keys
}

package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/edmartt/bookstatic-book-service/internal/adapters/http/dtos"
	"github.com/edmartt/bookstatic-book-service/internal/core/domain/models"
	"github.com/edmartt/bookstatic-book-service/internal/core/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ReadBook helpers

func newReadTestRouter(repository ports.IDataAccessLayer) *gin.Engine {
	handler := newTestHandler(repository)

	router := gin.New()
	router.GET("/api/v1/books/:id", handler.ReadBook)

	return router
}

func performReadRequest(router *gin.Engine, bookID, userID string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/books/"+bookID, nil)

	if userID != "" {
		request.Header.Set("X-User-Id", userID)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder
}

func sampleBook() *models.Books {
	return &models.Books{
		UUID:        "550e8400-e29b-41d4-a716-446655440000",
		ISBN:        "978-3-16-148410-0",
		Title:       "The Go Programming Language",
		Pages:       "380",
		CurrentPage: "120",
		Author:      "Donovan and Kernighan",
		Year:        "2015",
		Status:      "1",
		UserID:      "user-123",
	}
}

// ReadBook tests

func TestReadBookSuccess(t *testing.T) {
	book := sampleBook()
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return book, nil
		},
	}

	router := newReadTestRouter(repository)
	recorder := performReadRequest(router, book.UUID, book.UserID)

	require.Equal(t, nethttp.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Content-Type"), "application/json")

	var response dtos.BookResponseDTO
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))

	assert.Equal(t, dtos.BookResponseDTO{
		UUID:        book.UUID,
		ISBN:        book.ISBN,
		Title:       book.Title,
		Pages:       book.Pages,
		CurrentPage: book.CurrentPage,
		Author:      book.Author,
		Year:        book.Year,
		Status:      book.Status,
	}, response)
}

func TestReadBookDoesNotExposeUserID(t *testing.T) {
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return sampleBook(), nil
		},
	}

	router := newReadTestRouter(repository)
	recorder := performReadRequest(router, sampleBook().UUID, sampleBook().UserID)

	require.Equal(t, nethttp.StatusOK, recorder.Code)

	assert.ElementsMatch(t, []string{
		"id",
		"isbn",
		"title",
		"pages",
		"current_page",
		"author",
		"year",
		"status",
	}, responseKeys(t, recorder))
}

func TestReadBookPassesIDAndUserHeader(t *testing.T) {
	var receivedID string
	var receivedUserID string

	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			receivedID = id
			receivedUserID = userID

			return sampleBook(), nil
		},
	}

	router := newReadTestRouter(repository)
	recorder := performReadRequest(router, "book-uuid-42", "user-456")

	require.Equal(t, nethttp.StatusOK, recorder.Code)
	assert.Equal(t, "book-uuid-42", receivedID)
	assert.Equal(t, "user-456", receivedUserID)
}

// TestReadBookRepositoryError documents current behaviour: the handler maps any
// repository error to 404, including infrastructure failures that would ideally
// surface as 500.
func TestReadBookRepositoryError(t *testing.T) {
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return nil, errors.New("record not found")
		},
	}

	router := newReadTestRouter(repository)
	recorder := performReadRequest(router, "missing-book", "user-123")

	require.Equal(t, nethttp.StatusNotFound, recorder.Code)
	assert.Equal(t, "book not found", errorMessage(t, recorder))
}

func TestReadBookEmptyID(t *testing.T) {
	repositoryReadCalled := false
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			repositoryReadCalled = true

			return nil, errors.New("should not be called")
		},
	}

	handler := newTestHandler(repository)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(nethttp.MethodGet, "/api/v1/books/", nil)

	handler.ReadBook(context)

	require.Equal(t, nethttp.StatusBadRequest, recorder.Code)
	assert.False(t, repositoryReadCalled)
	assert.Equal(t, "book ID is empty", errorMessage(t, recorder))
}

func TestReadBookWithoutUserHeader(t *testing.T) {
	// The X-User-Id header is injected by Traefik's forward-auth, which rejects
	// requests without it before they reach this service. This test documents
	// that the handler itself performs no validation and forwards an empty
	// user id to the repository, whose result determines the response.
	var receivedUserID string
	receivedUserIDSet := false

	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			receivedUserID = userID
			receivedUserIDSet = true

			return nil, errors.New("record not found")
		},
	}

	router := newReadTestRouter(repository)
	recorder := performReadRequest(router, sampleBook().UUID, "")

	require.Equal(t, nethttp.StatusNotFound, recorder.Code)
	assert.True(t, receivedUserIDSet)
	assert.Equal(t, "", receivedUserID)
}

// CreateBook helpers

func newCreateTestRouter(repository ports.IDataAccessLayer) *gin.Engine {
	handler := newTestHandler(repository)

	router := gin.New()
	router.POST("/api/v1/books", handler.CreateBook)

	return router
}

func performCreateRequest(router *gin.Engine, body, userID string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/books", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	if userID != "" {
		request.Header.Set("X-User-Id", userID)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder
}

func validCreateFields() map[string]any {
	return map[string]any{
		"isbn":         "978-3-16-148410-0",
		"title":        "The Go Programming Language",
		"pages":        "380",
		"current_page": "120",
		"author":       "Donovan and Kernighan",
		"year":         "2015",
		"status":       "1",
	}
}

func validCreatePayload(t *testing.T) string {
	t.Helper()

	return marshalPayload(t, validCreateFields())
}

// CreateBook tests

func TestCreateBookSuccess(t *testing.T) {
	generatedID := uuid.NewString()
	repositoryCalled := false
	repository := fakeBookRepository{
		createFn: func(book models.Books) (string, error) {
			repositoryCalled = true

			return generatedID, nil
		},
	}

	router := newCreateTestRouter(repository)
	recorder := performCreateRequest(router, validCreatePayload(t), "user-123")

	require.Equal(t, nethttp.StatusCreated, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Content-Type"), "application/json")
	assert.True(t, repositoryCalled)

	var response BookCreatedResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, BookCreatedResponse{ID: generatedID}, response)
}

func TestCreateBookPassesPayloadToRepository(t *testing.T) {
	var received models.Books
	repository := fakeBookRepository{
		createFn: func(book models.Books) (string, error) {
			received = book

			return book.UUID, nil
		},
	}

	router := newCreateTestRouter(repository)
	recorder := performCreateRequest(router, validCreatePayload(t), "user-123")

	require.Equal(t, nethttp.StatusCreated, recorder.Code)

	assert.Equal(t, "978-3-16-148410-0", received.ISBN)
	assert.Equal(t, "The Go Programming Language", received.Title)
	assert.Equal(t, "380", received.Pages)
	assert.Equal(t, "120", received.CurrentPage)
	assert.Equal(t, "Donovan and Kernighan", received.Author)
	assert.Equal(t, "2015", received.Year)
	assert.Equal(t, "1", received.Status)
}

func TestCreateBookGeneratesUUID(t *testing.T) {
	var generatedIDs []string
	repository := fakeBookRepository{
		createFn: func(book models.Books) (string, error) {
			generatedIDs = append(generatedIDs, book.UUID)

			return book.UUID, nil
		},
	}

	router := newCreateTestRouter(repository)

	firstRecorder := performCreateRequest(router, validCreatePayload(t), "user-123")
	secondRecorder := performCreateRequest(router, validCreatePayload(t), "user-123")

	require.Equal(t, nethttp.StatusCreated, firstRecorder.Code)
	require.Equal(t, nethttp.StatusCreated, secondRecorder.Code)
	require.Len(t, generatedIDs, 2)

	for _, id := range generatedIDs {
		_, err := uuid.Parse(id)
		require.NoError(t, err, "generated id %q is not a valid UUID", id)
	}

	assert.NotEqual(t, generatedIDs[0], generatedIDs[1])
}

func TestCreateBookAssignsUserIDFromHeader(t *testing.T) {
	fields := validCreateFields()
	fields["user_id"] = "user-from-body"

	var received models.Books
	repository := fakeBookRepository{
		createFn: func(book models.Books) (string, error) {
			received = book

			return book.UUID, nil
		},
	}

	router := newCreateTestRouter(repository)
	recorder := performCreateRequest(router, marshalPayload(t, fields), "user-from-header")

	require.Equal(t, nethttp.StatusCreated, recorder.Code)
	assert.Equal(t, "user-from-header", received.UserID)
}

func TestCreateBookIgnoresClientSuppliedID(t *testing.T) {
	fields := validCreateFields()
	fields["id"] = "client-provided-id"

	var received models.Books
	repository := fakeBookRepository{
		createFn: func(book models.Books) (string, error) {
			received = book

			return book.UUID, nil
		},
	}

	router := newCreateTestRouter(repository)
	recorder := performCreateRequest(router, marshalPayload(t, fields), "user-123")

	require.Equal(t, nethttp.StatusCreated, recorder.Code)
	assert.NotEqual(t, "client-provided-id", received.UUID)

	_, err := uuid.Parse(received.UUID)
	assert.NoError(t, err)
}

func TestCreateBookInvalidJSON(t *testing.T) {
	repositoryCreateCalled := false
	repository := fakeBookRepository{
		createFn: func(book models.Books) (string, error) {
			repositoryCreateCalled = true

			return "", errors.New("should not be called")
		},
	}

	router := newCreateTestRouter(repository)
	recorder := performCreateRequest(router, `{"title":`, "user-123")

	require.Equal(t, nethttp.StatusBadRequest, recorder.Code)
	assert.False(t, repositoryCreateCalled)
	assert.Equal(t, "invalid request body", errorMessage(t, recorder))
}

func TestCreateBookEmptyBody(t *testing.T) {
	repositoryCreateCalled := false
	repository := fakeBookRepository{
		createFn: func(book models.Books) (string, error) {
			repositoryCreateCalled = true

			return "", errors.New("should not be called")
		},
	}

	router := newCreateTestRouter(repository)
	recorder := performCreateRequest(router, "", "user-123")

	require.Equal(t, nethttp.StatusBadRequest, recorder.Code)
	assert.False(t, repositoryCreateCalled)
	assert.Equal(t, "invalid request body", errorMessage(t, recorder))
}

func TestCreateBookMissingFields(t *testing.T) {
	testCases := []struct {
		name        string
		missingKeys []string
	}{
		{name: "all fields missing", missingKeys: []string{"isbn", "title", "pages", "current_page", "author", "year", "status"}},
		{name: "isbn missing", missingKeys: []string{"isbn"}},
		{name: "title missing", missingKeys: []string{"title"}},
		{name: "pages missing", missingKeys: []string{"pages"}},
		{name: "current_page missing", missingKeys: []string{"current_page"}},
		{name: "author missing", missingKeys: []string{"author"}},
		{name: "year missing", missingKeys: []string{"year"}},
		{name: "status missing", missingKeys: []string{"status"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fields := validCreateFields()

			for _, key := range testCase.missingKeys {
				delete(fields, key)
			}

			repositoryCreateCalled := false
			repository := fakeBookRepository{
				createFn: func(book models.Books) (string, error) {
					repositoryCreateCalled = true

					return "", errors.New("should not be called")
				},
			}

			router := newCreateTestRouter(repository)
			recorder := performCreateRequest(router, marshalPayload(t, fields), "user-123")

			require.Equal(t, nethttp.StatusBadRequest, recorder.Code)
			assert.False(t, repositoryCreateCalled)
			assert.Equal(t, "invalid request body", errorMessage(t, recorder))
		})
	}
}

func TestCreateBookNullField(t *testing.T) {
	fields := validCreateFields()
	fields["title"] = nil

	repositoryCreateCalled := false
	repository := fakeBookRepository{
		createFn: func(book models.Books) (string, error) {
			repositoryCreateCalled = true

			return "", errors.New("should not be called")
		},
	}

	router := newCreateTestRouter(repository)
	recorder := performCreateRequest(router, marshalPayload(t, fields), "user-123")

	require.Equal(t, nethttp.StatusBadRequest, recorder.Code)
	assert.False(t, repositoryCreateCalled)
	assert.Equal(t, "invalid request body", errorMessage(t, recorder))
}

func TestCreateBookMissingUserHeader(t *testing.T) {
	repositoryCreateCalled := false
	repository := fakeBookRepository{
		createFn: func(book models.Books) (string, error) {
			repositoryCreateCalled = true

			return "", errors.New("should not be called")
		},
	}

	router := newCreateTestRouter(repository)
	recorder := performCreateRequest(router, validCreatePayload(t), "")

	require.Equal(t, nethttp.StatusInternalServerError, recorder.Code)
	assert.False(t, repositoryCreateCalled)
	assert.Equal(t, "internal server error", errorMessage(t, recorder))
}

func TestCreateBookRepositoryError(t *testing.T) {
	repository := fakeBookRepository{
		createFn: func(book models.Books) (string, error) {
			return "", errors.New("insert failed")
		},
	}

	router := newCreateTestRouter(repository)
	recorder := performCreateRequest(router, validCreatePayload(t), "user-123")

	require.Equal(t, nethttp.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "error creating book", errorMessage(t, recorder))
}

func TestCreateBookResponseOnlyContainsID(t *testing.T) {
	repository := fakeBookRepository{
		createFn: func(book models.Books) (string, error) {
			return book.UUID, nil
		},
	}

	router := newCreateTestRouter(repository)
	recorder := performCreateRequest(router, validCreatePayload(t), "user-123")

	require.Equal(t, nethttp.StatusCreated, recorder.Code)

	assert.ElementsMatch(t, []string{"id"}, responseKeys(t, recorder))
}

// UpdateBook helpers

func newUpdateTestRouter(repository ports.IDataAccessLayer) *gin.Engine {
	handler := newTestHandler(repository)

	router := gin.New()
	router.PATCH("/api/v1/books/:id", handler.UpdateBook)

	return router
}

func performUpdateRequest(router *gin.Engine, bookID, body, userID string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(nethttp.MethodPatch, "/api/v1/books/"+bookID, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	if userID != "" {
		request.Header.Set("X-User-Id", userID)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder
}

func validUpdateFields() map[string]any {
	return map[string]any{
		"isbn":         "978-0-13-419044-0",
		"title":        "The Go Programming Language, 2nd edition",
		"pages":        "400",
		"current_page": "200",
		"author":       "Donovan and Kernighan",
		"year":         "2016",
		"status":       "2",
	}
}

func validUpdatePayload(t *testing.T) string {
	t.Helper()

	return marshalPayload(t, validUpdateFields())
}

// UpdateBook tests

func TestUpdateBookSuccess(t *testing.T) {
	var received models.Books
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return sampleBook(), nil
		},
		updateFn: func(book *models.Books) (*models.Books, error) {
			received = *book

			return book, nil
		},
	}

	router := newUpdateTestRouter(repository)
	recorder := performUpdateRequest(router, sampleBook().UUID, validUpdatePayload(t), sampleBook().UserID)

	require.Equal(t, nethttp.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Content-Type"), "application/json")
	assert.Equal(t, "The Go Programming Language, 2nd edition", received.Title)

	var response dtos.UpdateBookDTO
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))

	assert.Equal(t, dtos.UpdateBookDTO{
		ID:          stringPtr(sampleBook().UUID),
		ISBN:        stringPtr("978-0-13-419044-0"),
		Title:       stringPtr("The Go Programming Language, 2nd edition"),
		Pages:       stringPtr("400"),
		CurrentPage: stringPtr("200"),
		Author:      stringPtr("Donovan and Kernighan"),
		Year:        stringPtr("2016"),
		Status:      stringPtr("2"),
	}, response)
}

func TestUpdateBookAppliesPartialUpdate(t *testing.T) {
	testCases := []struct {
		name     string
		payload  map[string]any
		expected func(book *models.Books)
	}{
		{
			name:     "isbn",
			payload:  map[string]any{"isbn": "978-0-13-419044-0"},
			expected: func(book *models.Books) { book.ISBN = "978-0-13-419044-0" },
		},
		{
			name:     "title",
			payload:  map[string]any{"title": "Updated title"},
			expected: func(book *models.Books) { book.Title = "Updated title" },
		},
		{
			name:     "pages",
			payload:  map[string]any{"pages": "500"},
			expected: func(book *models.Books) { book.Pages = "500" },
		},
		{
			name:     "current_page",
			payload:  map[string]any{"current_page": "300"},
			expected: func(book *models.Books) { book.CurrentPage = "300" },
		},
		{
			name:     "author",
			payload:  map[string]any{"author": "Updated author"},
			expected: func(book *models.Books) { book.Author = "Updated author" },
		},
		{
			name:     "year",
			payload:  map[string]any{"year": "2020"},
			expected: func(book *models.Books) { book.Year = "2020" },
		},
		{
			name:     "status",
			payload:  map[string]any{"status": "3"},
			expected: func(book *models.Books) { book.Status = "3" },
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var received models.Books
			repository := fakeBookRepository{
				readFn: func(id, userID string) (*models.Books, error) {
					return sampleBook(), nil
				},
				updateFn: func(book *models.Books) (*models.Books, error) {
					received = *book

					return book, nil
				},
			}

			expected := sampleBook()
			testCase.expected(expected)

			router := newUpdateTestRouter(repository)
			recorder := performUpdateRequest(router, sampleBook().UUID, marshalPayload(t, testCase.payload), sampleBook().UserID)

			require.Equal(t, nethttp.StatusOK, recorder.Code)
			assert.Equal(t, *expected, received)
		})
	}
}

func TestUpdateBookNullFieldIsIgnored(t *testing.T) {
	var received models.Books
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return sampleBook(), nil
		},
		updateFn: func(book *models.Books) (*models.Books, error) {
			received = *book

			return book, nil
		},
	}

	router := newUpdateTestRouter(repository)
	recorder := performUpdateRequest(router, sampleBook().UUID, `{"title": null}`, sampleBook().UserID)

	require.Equal(t, nethttp.StatusOK, recorder.Code)
	assert.Equal(t, *sampleBook(), received)
}

func TestUpdateBookEmptyObject(t *testing.T) {
	var received models.Books
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return sampleBook(), nil
		},
		updateFn: func(book *models.Books) (*models.Books, error) {
			received = *book

			return book, nil
		},
	}

	router := newUpdateTestRouter(repository)
	recorder := performUpdateRequest(router, sampleBook().UUID, `{}`, sampleBook().UserID)

	require.Equal(t, nethttp.StatusOK, recorder.Code)
	assert.Equal(t, *sampleBook(), received)
}

func TestUpdateBookPassesIDAndUserHeader(t *testing.T) {
	var receivedID string
	var receivedUserID string

	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			receivedID = id
			receivedUserID = userID

			return sampleBook(), nil
		},
		updateFn: func(book *models.Books) (*models.Books, error) {
			return book, nil
		},
	}

	router := newUpdateTestRouter(repository)
	recorder := performUpdateRequest(router, "book-uuid-42", validUpdatePayload(t), "user-456")

	require.Equal(t, nethttp.StatusOK, recorder.Code)
	assert.Equal(t, "book-uuid-42", receivedID)
	assert.Equal(t, "user-456", receivedUserID)
}

func TestUpdateBookDoesNotExposeUserID(t *testing.T) {
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return sampleBook(), nil
		},
		updateFn: func(book *models.Books) (*models.Books, error) {
			return book, nil
		},
	}

	router := newUpdateTestRouter(repository)
	recorder := performUpdateRequest(router, sampleBook().UUID, validUpdatePayload(t), sampleBook().UserID)

	require.Equal(t, nethttp.StatusOK, recorder.Code)

	assert.ElementsMatch(t, []string{
		"id",
		"isbn",
		"title",
		"pages",
		"current_page",
		"author",
		"year",
		"status",
	}, responseKeys(t, recorder))
}

func TestUpdateBookReadError(t *testing.T) {
	// updateFn is intentionally left unset: the fake panics if Update is reached.
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return nil, errors.New("record not found")
		},
	}

	router := newUpdateTestRouter(repository)
	recorder := performUpdateRequest(router, "missing-book", validUpdatePayload(t), "user-123")

	require.Equal(t, nethttp.StatusNotFound, recorder.Code)
	assert.Equal(t, "book not found", errorMessage(t, recorder))
}

func TestUpdateBookNilBookFromRead(t *testing.T) {
	// Read is documented to return an error when the book is missing, but a
	// repository returning (nil, nil) must not panic the handler.
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return nil, nil
		},
	}

	router := newUpdateTestRouter(repository)
	recorder := performUpdateRequest(router, sampleBook().UUID, validUpdatePayload(t), "user-123")

	require.Equal(t, nethttp.StatusNotFound, recorder.Code)
	assert.Equal(t, "book not found", errorMessage(t, recorder))
}

func TestUpdateBookEmptyID(t *testing.T) {
	repositoryReadCalled := false
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			repositoryReadCalled = true

			return nil, errors.New("should not be called")
		},
	}

	handler := newTestHandler(repository)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(nethttp.MethodPatch, "/api/v1/books/", nil)

	handler.UpdateBook(context)

	require.Equal(t, nethttp.StatusBadRequest, recorder.Code)
	assert.False(t, repositoryReadCalled)
	assert.Equal(t, "id required", errorMessage(t, recorder))
}

func TestUpdateBookInvalidJSON(t *testing.T) {
	// The read happens before the bind, so readFn must be configured for the
	// handler to reach body validation. updateFn is left unset on purpose.
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return sampleBook(), nil
		},
	}

	router := newUpdateTestRouter(repository)
	recorder := performUpdateRequest(router, sampleBook().UUID, `{"title":`, "user-123")

	require.Equal(t, nethttp.StatusBadRequest, recorder.Code)
	assert.Equal(t, "invalid request body", errorMessage(t, recorder))
}

func TestUpdateBookEmptyBody(t *testing.T) {
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return sampleBook(), nil
		},
	}

	router := newUpdateTestRouter(repository)
	recorder := performUpdateRequest(router, sampleBook().UUID, "", "user-123")

	require.Equal(t, nethttp.StatusBadRequest, recorder.Code)
	assert.Equal(t, "invalid request body", errorMessage(t, recorder))
}

func TestUpdateBookRepositoryError(t *testing.T) {
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return sampleBook(), nil
		},
		updateFn: func(book *models.Books) (*models.Books, error) {
			return nil, errors.New("update failed")
		},
	}

	router := newUpdateTestRouter(repository)
	recorder := performUpdateRequest(router, sampleBook().UUID, validUpdatePayload(t), "user-123")

	require.Equal(t, nethttp.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "internal server error", errorMessage(t, recorder))
}

func TestUpdateBookNoRows(t *testing.T) {
	// documents current behaviour: a row deleted between Read and Update makes
	// the repository return sql.ErrNoRows, which the handler maps to 500.
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return sampleBook(), nil
		},
		updateFn: func(book *models.Books) (*models.Books, error) {
			return nil, sql.ErrNoRows
		},
	}

	router := newUpdateTestRouter(repository)
	recorder := performUpdateRequest(router, sampleBook().UUID, validUpdatePayload(t), "user-123")

	require.Equal(t, nethttp.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "internal server error", errorMessage(t, recorder))
}

func TestUpdateBookNilResult(t *testing.T) {
	// A repository returning (nil, nil) must not panic while building the response.
	repository := fakeBookRepository{
		readFn: func(id, userID string) (*models.Books, error) {
			return sampleBook(), nil
		},
		updateFn: func(book *models.Books) (*models.Books, error) {
			return nil, nil
		},
	}

	router := newUpdateTestRouter(repository)
	recorder := performUpdateRequest(router, sampleBook().UUID, validUpdatePayload(t), "user-123")

	require.Equal(t, nethttp.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "internal server error", errorMessage(t, recorder))
}

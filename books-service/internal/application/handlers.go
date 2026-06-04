package application

import (
	"net/http"

	"github.com/edmartt/bookstatic-book-service/internal/books/data"
	"github.com/edmartt/bookstatic-book-service/internal/books/dtos"
	"github.com/edmartt/bookstatic-book-service/internal/books/models"
	selfErrors "github.com/edmartt/booktastic-shared/errors/http/adapters/ginhttp"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HTTPHandler struct {
	bookRepository data.IDataAccessLayer
}

func NewHandler(bookRepo data.IDataAccessLayer) *HTTPHandler {
	return &HTTPHandler{
		bookRepository: bookRepo,
	}
}

func (h HTTPHandler) ReadBook(context *gin.Context) {
	id := context.Param("id")
	userID := context.GetHeader("X-User-Id")
	writer := selfErrors.NewGinErrors(context)

	if id == "" {
		writer.WriteError(http.StatusBadRequest, "book ID is empty")
		return
	}

	book, err := h.bookRepository.Read(id, userID)

	if err != nil {
		writer.WriteError(http.StatusNotFound, "book not found")
		return
	}
	response := dtos.BookResponseDTO{
		UUID:        book.UUID,
		ISBN:        book.ISBN,
		Title:       book.Title,
		Pages:       book.Pages,
		CurrentPage: book.CurrentPage,
		Author:      book.Author,
		Year:        book.Year,
		Status:      book.Status,
	}

	context.JSON(http.StatusOK, response)
}

func (h HTTPHandler) CreateBook(context *gin.Context) {
	var createDTO dtos.CreateBookDTO
	writer := selfErrors.NewGinErrors(context)

	if err := context.BindJSON(&createDTO); err != nil {
		writer.WriteError(http.StatusBadRequest, "invalid request body")
		return
	}

	userID := context.GetHeader("X-User-Id")

	if userID == "" {
		writer.WriteError(http.StatusInternalServerError, "internal server error")
		return
	}

	book := models.Books{
		UUID:        uuid.NewString(),
		ISBN:        *createDTO.ISBN,
		Title:       *createDTO.Title,
		Pages:       *createDTO.Pages,
		CurrentPage: *createDTO.CurrentPage,
		Author:      *createDTO.Author,
		Year:        *createDTO.Year,
		Status:      *createDTO.Status,
		UserID:      userID,
	}

	dbResponse, err := h.bookRepository.Create(book)

	if err != nil {
		writer.WriteError(http.StatusInternalServerError, "error creating book")
		return
	}

	context.JSON(http.StatusCreated, dbResponse)

}

func (h HTTPHandler) UpdateBook(context *gin.Context) {

	id := context.Param("id")
	userID := context.GetHeader("X-User-Id")
	writer := selfErrors.NewGinErrors(context)

	if id == "" {
		writer.WriteError(http.StatusInternalServerError, "internal server error")
		return
	}

	dbResponse, err := h.bookRepository.Read(id, userID)

	if err != nil {
		writer.WriteError(http.StatusNotFound, "book not found")
		return
	}

	var dtoUpdate dtos.UpdateBookDTO

	if err := context.ShouldBindJSON(&dtoUpdate); err != nil {
		writer.WriteError(http.StatusBadRequest, "invalid request body")
		return
	}

	dtoUpdate.ApplyTo(dbResponse)

	updateResult, err := h.bookRepository.Update(dbResponse)

	if err != nil {
		writer.WriteError(http.StatusInternalServerError, "internal server error")
		return
	}

	context.JSON(http.StatusOK, updateResult)
}

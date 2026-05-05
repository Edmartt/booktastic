package application

import (
	"net/http"

	"github.com/edmartt/bookstatic-book-service/internal/books/data"
	"github.com/edmartt/bookstatic-book-service/internal/books/dtos"
	"github.com/edmartt/bookstatic-book-service/internal/books/models"
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

	if id == "" {
		context.JSON(http.StatusBadRequest, "bad request")
		return
	}

	book, err := h.bookRepository.Read(id)

	if err != nil {
		context.JSON(http.StatusNotFound, "not found")
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

	err := context.BindJSON(&createDTO)

	if err != nil {
		context.JSON(http.StatusBadRequest, "bad request")
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
	}

	dbResponse, err := h.bookRepository.Create(book)

	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	context.JSON(http.StatusCreated, dbResponse)

}

func (h HTTPHandler) UpdateBook(context *gin.Context) {

	id := context.Param("id")

	if id == "" {
		jsonResponse := "bad request"
		context.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	dbResponse, err := h.bookRepository.Read(id)

	if err != nil {
		context.JSON(http.StatusNotFound, "not found")
		return
	}

	var dtoUpdate dtos.UpdateBookDTO
	err = context.ShouldBindJSON(&dtoUpdate)

	if err != nil {
		context.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	dtoUpdate.ApplyTo(dbResponse)

	updateResult, err := h.bookRepository.Update(dbResponse)

	if err != nil {
		context.JSON(http.StatusBadRequest, "bad request")
		return
	}

	context.JSON(http.StatusOK, updateResult)
}

package application

import (
	"log"
	"net/http"

	"github.com/edmartt/bookstatic-book-service/internal/books/data"
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
	response := httpResponse{}

	if id == "" {
		response.Response = "bad request"
		context.JSON(http.StatusBadRequest, response)
	}

	book, err := h.bookRepository.Read(id)

	if err != nil {
		response.Response = "not found"
		context.JSON(http.StatusNotFound, response)
		return
	}

	context.JSON(http.StatusOK, book)
}

func (h HTTPHandler) CreateBook(context *gin.Context) {
	book := models.Books{}
	book.UUID = uuid.NewString()

	err := context.BindJSON(&book)
	jsonResponse := httpResponse{}

	if err != nil {
		jsonResponse.Response = "bad request"
		context.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	dbResponse, err := h.bookRepository.Create(book)

	if err != nil {
		log.Println(err.Error())
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	context.JSON(http.StatusCreated, dbResponse)

}

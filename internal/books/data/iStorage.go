package data

import "github.com/edmartt/bookstatic-book-service/internal/books/models"

type IDataAccessLayer interface {
	Create(book models.Books) (string, error)
	Read(id string) (*models.Books, error)
	Update(book *models.Books) (*models.Books, error)
}

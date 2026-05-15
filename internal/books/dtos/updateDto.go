package dtos

import "github.com/edmartt/bookstatic-book-service/internal/books/models"

type UpdateBookDTO struct {
	ISBN        *string `json:"isbn"`
	Title       *string `json:"title"`
	Pages       *string `json:"pages"`
	CurrentPage *string `json:"current_page"`
	Author      *string `json:"author"`
	Year        *string `json:"year"`
	Status      *string `json:"status"`
}

func (u UpdateBookDTO) ApplyTo(b *models.Books) {

	if u.ISBN != nil {
		b.ISBN = *u.ISBN
	}
	if u.Title != nil {
		b.Title = *u.Title
	}
	if u.Pages != nil {
		b.Pages = *u.Pages
	}
	if u.CurrentPage != nil {
		b.CurrentPage = *u.CurrentPage
	}
	if u.Author != nil {
		b.Author = *u.Author
	}
	if u.Year != nil {
		b.Year = *u.Year
	}
	if u.Status != nil {
		b.Status = *u.Status
	}

}

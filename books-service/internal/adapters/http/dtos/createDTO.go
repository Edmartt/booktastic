package dtos

type CreateBookDTO struct {
	UUID        string  `json:"id"`
	ISBN        *string `json:"isbn"`
	Title       *string `json:"title"`
	Pages       *string `json:"pages"`
	CurrentPage *string `json:"current_page"`
	Author      *string `json:"author"`
	Year        *string `json:"year"`
	Status      *string `json:"status"`
}

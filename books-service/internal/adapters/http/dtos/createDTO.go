package dtos

type CreateBookDTO struct {
	UUID        string  `json:"id"`
	ISBN        *string `json:"isbn" binding:"required"`
	Title       *string `json:"title" binding:"required"`
	Pages       *string `json:"pages" binding:"required"`
	CurrentPage *string `json:"current_page" binding:"required"`
	Author      *string `json:"author" binding:"required"`
	Year        *string `json:"year" binding:"required"`
	Status      *string `json:"status" binding:"required"`
}

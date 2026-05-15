package models

// I want to know if I'm reading, I read or if I need to buy that book, Status only allowed answers
const (
	Read = iota
	Reading
	ToBeRead
	PendingPurchase
)

// Books model
type Books struct {
	UUID        string `db:"uuid"`
	ISBN        string `db:"isbn"`
	Title       string `db:"title"`
	Pages       string `db:"pages"`
	CurrentPage string `db:"current_page"`
	Author      string `db:"author"`
	Year        string `db:"year"`
	Status      string `db:"status"`
}

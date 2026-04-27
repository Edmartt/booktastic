package data

import (
	"fmt"
	"log"

	"github.com/edmartt/bookstatic-book-service/internal/books/models"
	"github.com/edmartt/bookstatic-book-service/internal/database"
)

type BookDataAccess struct {
	db   database.IDBConnection
	book models.Books
}

func NewRepository(db database.IDBConnection) *BookDataAccess {
	return &BookDataAccess{
		db: db,
	}
}

func (b BookDataAccess) Create(book models.Books) (string, error) {
	conn := b.db.GetConnection()

	_, err := conn.NamedExec("INSERT INTO books (uuid, isbn, title, pages, current_page, author, year, status) VALUES(:uuid, :isbn, :title, :pages, :current_page, :author, :year, :status)", &book)

	if err != nil {
		return "", err
	}

	bookData := fmt.Sprintf("%s %s %s", book.UUID, book.ISBN, book.Title)

	return bookData, nil
}

func (b BookDataAccess) Read(id string) (*models.Books, error) {
	conn := b.db.GetConnection()

	query := "SELECT uuid, isbn, title, pages, current_page, author, year, status FROM books WHERE uuid = ?"

	query = conn.Rebind(query)

	err := conn.Get(&b.book, query, id)

	if err != nil {
		log.Println("error data: ", err.Error())
		return nil, err
	}
	return &b.book, nil
}

func (b BookDataAccess) Update(id string) (string, error) {
	conn := b.db.GetConnection()
	book, err := b.Read(id)

	if err != nil {
		return "", err
	}

	result, err := conn.NamedExec("UPDATE books SET isbn = :isbn, title = :title, pages = :pages, current_page = :current_page, author = :author, year = :year, status = :status WHERE id=?", book)

	if err != nil {
		return "", err
	}

	lastID, err := result.LastInsertId()

	if err != nil {
		return "", err
	}

	bookData := fmt.Sprintf("%d %s %s", lastID, book.ISBN, book.Title)

	return bookData, nil
}

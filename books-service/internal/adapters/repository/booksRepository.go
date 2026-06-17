package repository

import (
	"database/sql"
	"log"

	"github.com/edmartt/bookstatic-book-service/internal/core/domain/models"
	"github.com/edmartt/bookstatic-book-service/internal/core/ports"
)

type BookDataAccess struct {
	db ports.IDBConnection
}

func NewRepository(db ports.IDBConnection) *BookDataAccess {
	return &BookDataAccess{
		db: db,
	}
}

func (b *BookDataAccess) Create(book models.Books) (string, error) {
	conn := b.db.GetConnection()

	_, err := conn.NamedExec("INSERT INTO books (uuid, isbn, title, pages, current_page, author, year, status, user_id) VALUES(:uuid, :isbn, :title, :pages, :current_page, :author, :year, :status, :user_id)", &book)

	if err != nil {
		return "", err
	}

	return book.UUID, nil
}

func (b *BookDataAccess) Read(id, user_id string) (*models.Books, error) {
	conn := b.db.GetConnection()

	query := "SELECT uuid, isbn, title, pages, current_page, author, year, status FROM books WHERE uuid = ? AND user_id = ?"

	query = conn.Rebind(query)

	var book models.Books

	err := conn.Get(&book, query, id, user_id)

	if err != nil {
		log.Println("error data: ", err.Error())
		return nil, err
	}
	return &book, nil
}

func (b *BookDataAccess) Update(book *models.Books) (*models.Books, error) {
	conn := b.db.GetConnection()

	result, err := conn.NamedExec("UPDATE books SET isbn = :isbn, title = :title, pages = :pages, current_page = :current_page, author = :author, year = :year, status = :status WHERE uuid = :uuid", book)

	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, sql.ErrNoRows
	}

	return book, nil
}

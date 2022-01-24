package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	userTable                = "users"
	bookTable                = "books"
	bookCopiesTable          = "book_copies"
	bookBookCopiesTable      = "books_book_copies"
	bookPhotoTable           = "books_photo"
	bookBookPhotoTable       = "books_books_photo"
	authorsTable             = "authors"
	bookAuthorsTable         = "books_authors"
	authorsPhotoTable        = "authors_photo"
	authorsAuthorsPhotoTable = "authors_author_photo"
	genreTable               = "genres"
	bookGenreTable           = "book_genre"
	bookUserTable            = "books_users"
	limit                    = 40
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

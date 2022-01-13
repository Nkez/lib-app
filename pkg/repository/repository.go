package repository

import (
	library_app "github.com/Nkez/lib-app.git"
	"github.com/jmoiron/sqlx"
)

type Cart interface {
	CreateCart(user library_app.OrderCart) (library_app.OrderCart, error)

	CheckBooks(book string) (checkBook string, err error)
	GetPrice(book string) (priceBook float64, err error)
	GetUser(email string) (library_app.OrderCart, error)
	CheckOrderBook(email string) (orderBook []string, err error)
	FindRegisterUser(email string) (library_app.User, error)
	GetEmailToSend() ([]library_app.OrderCart, error)
	UpdatePrice()
}

type Book interface {
	CreateBook(book library_app.Book) (int, error)
}

type User interface {
	CreateUser(user library_app.User) (int, error)
	GetAllUsers() ([]library_app.User, error)
}

type ReturnBook interface {
	GetAllBooks() ([]library_app.ReturnBook, error)
}

type ReturnCart interface {
	CreateRtCart() error
}

type Repository struct {
	Book
	User
	Cart
	ReturnBook
	ReturnCart
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:       NewUserPostgres(db),
		Book:       NewBookPostgres(db),
		Cart:       NewCartPostgres(db),
		ReturnBook: NewReturnBookPostgres(db),
		ReturnCart: NewReturnCartPostgrers(db),
	}
}

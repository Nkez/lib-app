package services

import (
	library_app "github.com/Nkez/lib-app.git"
	"github.com/Nkez/lib-app.git/pkg/repository"
)

type Cart interface {
	CreateCart(email string, books []string) (library_app.OrderCart, error)
	CheckBooks(book string) (checkBook string, err error)
	GetPrice(book string) (priceBook float64, err error)
	GetUser(email string) (library_app.OrderCart, error)
	CheckOrderBook(email string) (orderBook []string, err error)
	FindRegisterUser(email string) (library_app.User, error)

	GetEmailToSend() ([]library_app.OrderCart, []string, error)
	UpdatePrice() error
	SendEmail()
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

type Service struct {
	Book
	User
	Cart
	ReturnBook
	ReturnCart
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		User:       NewUserService(repository.User),
		Book:       NewBookService(repository.Book),
		Cart:       NewCartService(repository.Cart),
		ReturnBook: NewReturnBookService(repository.ReturnBook),
		ReturnCart: NewReturnCartService(repository.ReturnCart),
	}
}

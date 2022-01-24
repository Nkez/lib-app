package services

import (
	"github.com/Nkez/library-app.git/models"
	"github.com/Nkez/library-app.git/pkg/repository"
)

type Book interface {
	CreateBook(book models.Book) (int, error)
	CreateGenre(genres models.CreateGenre) (id int, err error)
	CreateAuthor(authors models.CreateAuthor) (id int, err error)
	JoinBookPhoto(idBook int, p, n []string) error
	GetAllBooks(page, limit string) ([]models.ReturnBook, error)
	GetAllGenres() ([]string, error)
	GetByTitle(title string) ([]models.ReturnBook, error)
	GetByWord(word string) ([]string, error)

	GetAutPhoto(id int) (path, photo string, err error)

	GetAutInfo(id int) (models.CreateAuthor, error)
	ChangeAutPhoto(author models.CreateAuthor)
	ChangeBookPhoto(newPath, newName, paths string)

	BookPhoto(idPhoto int) (path, name string, err error)
	GetByTopRating() []models.TopRating
	JoinDefetBookPhoto(idBook int, defect string, paths []string) error
	GetDefectPhoto(idPhoto int) (string, error)
	ChangeAuthorPhoto(id int, paths string) error

	DeleteGenre(id int) error
	//ChekRegisterUser(email string) (models.User, error)
}

type User interface {
	CreateUser(user models.User) (int, error)
	GetAllUsers(page, limit string) ([]models.User, error)
	GetByName(name string) []models.UserName
	FindByWord(name string) []models.UserName
}

type Order interface {
	CreateOrder(input models.OrderInput) (models.Order, error)
	GetAllOrder(page, limit string) ([]models.InfoOrdDept, error)
}

type Return interface {
	ReturnCart(input models.ReturnInput) (models.DbrInfo, error)
}
type Debors interface {
	GetAllDebors() ([]models.InfoOrdDept, error)
	FirstCheck()
	WaitAndEmailAgain()
}

type Repository struct {
	Book
	User
	Order
	Return
	Debors
}

type Service struct {
	Book
	User
	Order
	Return
	Debors
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		User:   NewUserService(repository.User),
		Book:   NewBookService(repository.Book),
		Order:  NewOrderService(repository.Order),
		Return: NewReturnService(repository.Return),
		Debors: NewDeborsService(repository.Debors),
	}
}

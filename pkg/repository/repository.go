package repository

import (
	"github.com/Nkez/library-app.git/models"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

type Book interface {
	//Insert
	InsertBookInfo(yearOfPublished, RegistrationDate time.Time, book models.Book) (int, error)
	InsertGenre(genres models.CreateGenre) (int, error)
	InsertAuthorInfo(author models.CreateAuthor) (int, error)
	JoinBookPhoto(idBook int, p, n []string) error
	//Join
	JoinBookGenres(id int, book models.Book) error
	JoinBookAuthor(id int, book models.Book) error
	GetAllBooks(page, limit string) ([]models.ReturnBook, error)
	ChekRegisterUser(email string) (models.User, error)
	GetAllGenres() ([]string, error)
	GetByTitle(title string) (models.ReturnBook, error)
	GetByWord(word string) ([]string, error)

	GetAutPhoto(id int) (path, photo string, err error)

	BookPhoto(idPhoto int) (path, name string, err error)
	GetByTopRating() []models.TopRating
	JoinDefetBookPhoto(idBook int, defect string, paths []string) error
	GetDefectPhoto(idPhoto int) (string, error)
	ChangeAuthorPhoto(id int, paths string) error

	GetAutInfo(id int) (models.CreateAuthor, error)
	ChangeAutPhoto(author models.CreateAuthor)

	ChangeBookPhoto(newPath, newName, paths string)

	DeleteGenre(id int) error
	GetBookById(bookId int) (models.ReturnBook, error)
	GetAllCopies(page, limit string) ([]models.CopiesInfo, error)
	GetCopieByBookIdAndId(bookId, copId string) (models.CopiesSolo, error)
}

type User interface {
	CreateUser(user models.User) (int, error)
	GetAllUsers(page, limit string) ([]models.User, error)
	GetByName(name string) []models.UserName
	FindByWord(name string) []models.UserName
}

type Order interface {
	FindUser(email string) (models.User, error)
	GetPrice(id []int) (prices []float64, err error)
	CheckIsReturn(id int) (books []string, err error)
	GetBooksId(books []string) (idBooks []int, err error)
	JoinBookUser(idUser int, idBook []int, price float64, orderDate, returnDate time.Time) error
	MinusInventoryCount(idBooks []int) error
	ReturnOrder(id int) (models.ReturnOrder, error)
	GetAllOrder(page, limit string) ([]models.ReturnOrder, error)

	JoinBookCopies(idBook []int) (idCopies []int, err error)
}

type Return interface {
	GetRating(idBooks []int) ([]float64, error)
	ReturnBooksInLibUpdateRating(booksId []int, rating []decimal.Decimal) error
	ReturnBook(idUser int, idBook []int, returnDate time.Time) (models.DbrInfo, error)
	CheckReturningBook(id int) ([]int, error)
	UpdatePrice(id int, newPrice float64)
}

type Debors interface {
	GetAllDebors() ([]models.InfoOrdDept, error)
}

type Repository struct {
	Book
	User
	Order
	Return
	Debors
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:   NewUserPostgres(db),
		Book:   NewBookPostgres(db),
		Order:  NewOrderPostgres(db),
		Return: NewReturnPostgres(db),
		Debors: NewDeborsPostgres(db),
	}
}

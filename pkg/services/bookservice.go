package services

import (
	library_app "github.com/Nkez/lib-app.git"
	"github.com/Nkez/lib-app.git/pkg/repository"
	"github.com/sirupsen/logrus"
	"time"
)

type BookService struct {
	repository repository.Book
}

type ReturnBookService struct {
	repository repository.ReturnBook
}

func NewReturnBookService(repository repository.ReturnBook) *ReturnBookService {
	return &ReturnBookService{repository: repository}
}

func NewBookService(repository repository.Book) *BookService {
	return &BookService{repository: repository}
}

func (s *BookService) CreateBook(book library_app.Book) (int, error) {
	logrus.Info("Create Book service")
	time.Parse("2006-01-02", book.YearOfPublished)
	time.Parse("2006-01-02", book.RegistrationDate)

	return s.repository.CreateBook(book)
}

func (s *ReturnBookService) GetAllBooks() ([]library_app.ReturnBook, error) {
	logrus.Info("Find All Books service")

	return s.repository.GetAllBooks()
}

package services

import (
	"fmt"
	"github.com/Nkez/library-app.git/models"
	"github.com/Nkez/library-app.git/pkg/repository"
	"github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"
	"time"
)

type BookService struct {
	repositoryBook repository.Book
}

func NewBookService(repositoryBook repository.Book) *BookService {
	return &BookService{repositoryBook: repositoryBook}
}

///Book

func (s *BookService) CreateBook(book models.Book) (int, error) {
	logrus.Info("create book service")
	logrus.Info("validate struct book")
	res, err := govalidator.ValidateStruct(book)
	if err != nil {
		println("error: " + err.Error())
	}
	if res == false {
		return 0, err
	}
	logrus.Info("insert book info")
	id, err := s.InsertBookInfo(book)
	if err != nil {
		return 0, err
	}
	logrus.Info(fmt.Sprintf("join book:id %v", id))
	s.JoinBookGenres(id, book)
	logrus.Info(fmt.Sprintf("join book:id %v", id))
	s.JoinBookAuthor(id, book)

	return id, nil
}
func (s *BookService) GetAllBooks(page, limit string) ([]models.ReturnBook, error) {
	logrus.Info("find all book service")
	return s.repositoryBook.GetAllBooks(page, limit)
}
func (s *BookService) JoinBookPhoto(idBook int, p, n []string) error {
	logrus.Info(fmt.Sprintf("join book id %v", idBook))
	return s.repositoryBook.JoinBookPhoto(idBook, p, n)
}
func (s *BookService) JoinDefetBookPhoto(idBook int, defect string, paths []string) error {
	return s.repositoryBook.JoinDefetBookPhoto(idBook, defect, paths)
}
func (s *BookService) BookPhoto(idPhoto int) (path, name string, err error) {
	path, name, err = s.repositoryBook.BookPhoto(idPhoto)
	if err != nil {
		return "", "", err
	}
	return path, name, nil
}
func (s *BookService) ChangeBookPhoto(newPath, newName, paths string) {
	s.repositoryBook.ChangeBookPhoto(newPath, newName, paths)
}
func (s *BookService) GetDefectPhoto(idPhoto int) (string, error) {
	path, err := s.repositoryBook.GetDefectPhoto(idPhoto)
	if err != nil {
		return "", err
	}
	return path, nil
}
func (s *BookService) InsertBookInfo(book models.Book) (int, error) {
	logrus.Info("parse  date")
	yearPub, err := time.Parse("2006-01-02", book.YearOfPublished)
	if err != nil {
		return 0, fmt.Errorf("write date in format YYYY-MM-DD, not %s", book.YearOfPublished)
	}

	regDate, err := time.Parse("2006-01-02", book.RegistrationDate)
	if err != nil {
		return 0, fmt.Errorf("write date in format YYYY-MM-DD, not %s", book.RegistrationDate)
	}
	id, err := s.repositoryBook.InsertBookInfo(yearPub, regDate, book)
	if err != nil {
		return 0, err
	}

	return id, nil
}

///Author
func (s *BookService) GetAutPhoto(id int) (path, photo string, err error) {
	path, photo, err = s.repositoryBook.GetAutPhoto(id)
	if err != nil {
		return "", "", err
	}
	return path, photo, nil
}
func (s *BookService) CreateAuthor(authors models.CreateAuthor) (id int, err error) {
	return s.repositoryBook.InsertAuthorInfo(authors)
}
func (s *BookService) JoinBookAuthor(id int, book models.Book) error {
	return s.repositoryBook.JoinBookAuthor(id, book)
}
func (s *BookService) ChangeAuthorPhoto(id int, paths string) error {
	return s.repositoryBook.ChangeAuthorPhoto(id, paths)
}
func (s *BookService) GetAutInfo(id int) (models.CreateAuthor, error) {
	return s.repositoryBook.GetAutInfo(id)
}
func (s *BookService) ChangeAutPhoto(author models.CreateAuthor) {
	s.repositoryBook.ChangeAutPhoto(author)
}

//Genre

func (s *BookService) GetAllGenres() ([]string, error) {
	logrus.Info("Find genre serive")
	return s.repositoryBook.GetAllGenres()
}
func (s *BookService) CreateGenre(genres models.CreateGenre) (id int, err error) {
	return s.repositoryBook.InsertGenre(genres)
}
func (s *BookService) JoinBookGenres(id int, book models.Book) error {
	s.repositoryBook.JoinBookGenres(id, book)
	return nil
}
func (s *BookService) DeleteGenre(id int) error {
	return s.repositoryBook.DeleteGenre(id)
}

//Find smt..
func (s *BookService) GetByWord(word string) ([]string, error) {
	return s.repositoryBook.GetByWord(word)
}
func (s *BookService) GetByTitle(title string) (models.ReturnBook, error) {
	return s.repositoryBook.GetByTitle(title)
}
func (s *BookService) GetByTopRating() []models.TopRating {
	return s.repositoryBook.GetByTopRating()
}
func (s *BookService) GetBookById(bookId int) (models.ReturnBook, error) {
	return s.repositoryBook.GetBookById(bookId)

}

//Copies

func (s *BookService) GetAllCopies(page, limit string) ([]models.CopiesInfo, error) {
	return s.repositoryBook.GetAllCopies(page, limit)
}

func (s *BookService) GetCopieByBookIdAndId(bookId, copId string) (models.CopiesSolo, error) {
	return s.repositoryBook.GetCopieByBookIdAndId(bookId, copId)
}

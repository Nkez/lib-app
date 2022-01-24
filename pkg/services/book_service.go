package services

import (
	"fmt"
	"github.com/Nkez/lib-app.git/models"
	"github.com/Nkez/lib-app.git/pkg/repository"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type BookService struct {
	repositoryBook repository.Book
}

func NewBookService(repositoryBook repository.Book) *BookService {
	return &BookService{repositoryBook: repositoryBook}
}

func (s *BookService) CreateBook(book models.Book) (int, error) {
	logrus.Info("Start  Create Book service.................")
	logrus.Info("Insert Book Info..............")
	id, err := s.InsertBookInfo(book)
	if err != nil {
		return 0, err
	}
	logrus.Info("Insert BookFoto................")
	_, err = s.InertBookFoto(book)
	if err != nil {
		return 0, err
	}
	logrus.Info("Join Book-BookFoto................")
	s.JoinBookBookFoto(book)

	logrus.Info("Insert Authors................")
	_, err = s.InsertAuthorInfo(book)
	if err != nil {
		return 0, err
	}
	logrus.Info("Join Book-Authors................")
	s.JoinBookAuthors(book)
	if err != nil {
		return 0, err
	}
	logrus.Info("Join Book-Genre................")
	err = s.JoinBookGenre(book)
	if err != nil {
		return 0, err
	}
	return id, nil
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
		return 0, err /*fmt.Errorf("error in book info")*/
	}
	return id, nil
}

func (s *BookService) InertBookFoto(book models.Book) ([]int, error) {
	var idSlice []int
	var id int
	fmt.Println(book.BooksPhoto)
	for i := 0; i < len(book.BooksPhoto); i++ {
		book.BooksPhoto[i].BookPhoto = strings.Trim(book.BooksPhoto[i].BookPhoto, "[{ }}")
		id, _ = s.repositoryBook.InertBookFoto(i, book)
		idSlice = append(idSlice, id)
	}
	return idSlice, nil
}

func (s *BookService) JoinBookBookFoto(book models.Book) {
	idBook, _ := s.InsertBookInfo(book)
	idFotos, _ := s.InertBookFoto(book)
	fmt.Println(idBook)
	fmt.Println(idFotos)
	for i := 0; i < len(idFotos); i++ {
		s.repositoryBook.JoinBookBookFoto(idBook, idFotos[i])
	}
}

func (s *BookService) InsertAuthorInfo(book models.Book) ([]int, error) {
	fmt.Println(book.Authors)
	var id int
	var idSlice []int
	for i := 0; i < len(book.Authors); i++ {
		id, _ = s.repositoryBook.InsertAuthorInfo(i, book)
		idSlice = append(idSlice, id)
	}
	fmt.Println(idSlice)
	return idSlice, nil
}

func (s *BookService) JoinBookAuthors(book models.Book) {
	idBook, _ := s.InsertBookInfo(book)
	idAuthors, _ := s.InsertAuthorInfo(book)
	for i := 0; i < len(idAuthors); i++ {
		s.repositoryBook.JoinBookAuthors(idBook, idAuthors[i])
	}
}

func (s *BookService) JoinBookGenre(book models.Book) error {
	idBook, _ := s.InsertBookInfo(book)
	var idGenres []int
	for i := 0; i < len(book.Genres); i++ {
		id, err := s.ParseAndCheckGenre(i, book)
		if err != nil {
			return err
		}
		idGenres = append(idGenres, id)
	}

	for i := 0; i < len(idGenres); i++ {
		s.repositoryBook.JoinBookGenre(idBook, idGenres[i])
	}
	return nil
}

func (s *BookService) ParseAndCheckGenre(index int, book models.Book) (id int, err error) {
	logrus.Info("parse genres................")
	book.Genres[index].Genre = strings.ToLower(strings.Trim(book.Genres[index].Genre, "{ }"))
	switch book.Genres[index].Genre {
	case "fantasy":
		return 1, nil
	case "adventure":
		return 2, nil
	case "romance":
		return 3, nil
	case "contemporary":
		return 4, nil
	case "dystopian":
		return 5, nil
	case "mystery":
		return 6, nil
	case "horror":
		return 7, nil
	case "thriller":
		return 8, nil
	case "paranormal":
		return 9, nil
	default:
		return 0, fmt.Errorf("choose a genre from the list")
	}
}

func (s *BookService) GetAllBooks() ([]models.ReturnBook, error) {
	logrus.Info("Find All Books service")
	return s.repositoryBook.GetAllBooks()

}

func (s *BookService) ChekRegisterUser(email string) (models.User, error) {

	return s.repositoryBook.ChekRegisterUser(email)
}


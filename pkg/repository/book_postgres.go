package repository

import (
	"fmt"
	library_app "github.com/Nkez/lib-app.git"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type BookPostgres struct {
	db *sqlx.DB
}

type ReturnBookPostgres struct {
	db *sqlx.DB
}

func NewReturnBookPostgres(db *sqlx.DB) *ReturnBookPostgres {
	return &ReturnBookPostgres{db: db}
}

func NewBookPostgres(db *sqlx.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

func (r *BookPostgres) CreateBook(book library_app.Book) (int, error) {
	logrus.Info("Create Book DB")
	var id int

	b := struct {
		Id               int
		BookTitle        string
		BookTitleNative  string
		BookPrice        decimal.Decimal
		InventoryCount   int
		TotalCount       int
		OneDayPrice      decimal.Decimal
		YearOfPublished  string
		RegistrationDate string
		NumberOfPages    int
		BookState        bool
		HideBook         bool
		Rating           int
		BookFoto         string
	}{
		Id:               book.Id,
		BookTitle:        book.BookTitle,
		BookTitleNative:  book.BookTitleNative,
		BookPrice:        book.BookPrice,
		InventoryCount:   book.InventoryCount,
		TotalCount:       book.TotalCount,
		OneDayPrice:      book.OneDayPrice,
		YearOfPublished:  book.YearOfPublished,
		RegistrationDate: book.RegistrationDate,
		NumberOfPages:    book.NumberOfPages,
		BookState:        book.HideBook,
		HideBook:         book.HideBook,
		Rating:           book.Rating,
		BookFoto:         book.BookFoto,
	}
	genre := struct {
		ActionAndAdventureGenre string
		//ClassicGenre             string
		//DetectiveAndMysteryGenre string
		//FantasyGenre             string
		//HistoricalFictionGenre   string
		//HorrorGenre              string
	}{
		ActionAndAdventureGenre: book.Genres.ActionAndAdventure,
		//ClassicGenre:             book.Genres.Classic,
		//DetectiveAndMysteryGenre: book.Genres.DetectiveAndMystery,
		//FantasyGenre:             book.Genres.Fantasy,
		//HistoricalFictionGenre:   book.Genres.HistoricalFiction,
		//HorrorGenre:              book.Genres.Horror,
	}

	query := fmt.Sprintf(`
				INSERT INTO books(book_title,book_title_native,book_price,
				inventory_count,books_in_lib,one_day_price,year_of_published,registration_date,
				number_of_pages,book_state,hide_book,rating,
				book_foto) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`)
	row := r.db.QueryRow(query, b.BookTitle, b.BookTitleNative, b.BookPrice,
		b.InventoryCount, b.TotalCount, b.OneDayPrice, b.YearOfPublished /*book.YearOfPublished*/, b.RegistrationDate, /*book.RegistrationDate*/
		b.NumberOfPages, b.BookState, b.HideBook, b.Rating, b.BookFoto)
	if err := row.Scan(&id); err != nil {
		return id, err
	}
	query = fmt.Sprintf(`INSERT INTO genres(genre) VALUES($1) RETURNING id`)
	row = r.db.QueryRow(query, genre.ActionAndAdventureGenre)
	if err := row.Scan(&id); err != nil {
		return id, err
	}
	for _, v := range book.Authors {
		query = fmt.Sprintf(`INSERT INTO authors(authors_firstname,authors_lastname,authors_foto) VALUES($1,$2,$3) RETURNING id`)
		row = r.db.QueryRow(query, v.FirstName, v.LastName, v.AutorFoto)
		if err := row.Scan(&id); err != nil {
			return id, err
		}
	}

	return id, nil
}

func (r *ReturnBookPostgres) GetAllBooks() ([]library_app.ReturnBook, error) {
	logrus.Info("Find All Books DB")
	var books []library_app.ReturnBook

	query := fmt.Sprintf(`SELECT  b.book_title, g.genre ,b.year_of_published, b.inventory_count,b.books_in_lib
									FROM genres g LEFT JOIN books b ON g.id = b.id 
									ORDER BY inventory_count DESC`)
	if err := r.db.Select(&books, query); err != nil {
		return books, err
	}
	return books, nil
}

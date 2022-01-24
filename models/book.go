package models

import (
	"github.com/shopspring/decimal"
)

type Book struct {
	Id               int             `json:"id" db:"id" valid:"required,numeric"`
	BookTitle        string          `json:"book_title" db:"book_title" valid:"required"`
	BookTitleNative  string          `json:"book_title_native" db:"book_title_native"`
	BookPrice        decimal.Decimal `json:"book_price" db:"book_price"`
	InventoryCount   int             `json:"inventory_count" db:"inventory_count"valid:"required,numeric"`
	TotalCount       int             `json:"books_in_lib" db:"books_in_lib"valid:"required,numeric"`
	OneDayPrice      decimal.Decimal `json:"one_day_price" db:"one_day_price"`
	YearOfPublished  string          `json:"year_of_published"  db:"year_of_published"`
	RegistrationDate string          `json:"registration_date"  db:"registration_date"`
	NumberOfPages    int             `json:"number_of_pages" db:"number_of_pages"`
	BookState        bool            `json:"book_state" db:"book_state"`
	HideBook         bool            `json:"hide_book" db:"hide_book"`
	Rating           decimal.Decimal `json:"rating" db:"rating"`
	Authors          []Authors       `json:"Authors"`
	BooksPhoto       []BooksPhoto    `json:"BooksPhoto"`
	Genres           []Genres        `json:"Genres"`
}
type Authors struct {
	Id int `json:"authors_id"`
}

type Genres struct {
	Genre int `json:"genre" db:"genre"`
}

type BooksPhoto struct {
	IdPhoto int `json:"id_photo" db:"id_photo"`
}

type CreateAuthor struct {
	Id          int    `json:"authors_id" db:"id"`
	FirstName   string `json:"authors_firstname" db:"author_firstname" valid:"required"`
	LastName    string `json:"authors_lastname" db:"author_lastname" valid:"required"`
	PhotoName   string `json:"photo_name" db:"foto_name"`
	AuthorPhoto string `json:"authors_photo" db:"authors_photo"`
}
type CreateGenre struct {
	Genre string `json:"genre" db:"genre"`
}

type ReturnBook struct {
	Id              int
	BookTitle       string   `json:"book_title" db:"book_title"`
	YearOfPublished string   `json:"year_of_published"  db:"year_of_published"`
	Genre           []string `json:"genre" db:"genre"`
	TotalCount      int      `json:"books_in_lib" db:"books_in_lib"`
	InventoryCount  int      `json:"inventory_count" db:"inventory_count"`
}

type ByWord struct {
	Book string
}

type TopRating struct {
	Id              int
	BookTitle       string          `json:"book_title" db:"book_title"`
	YearOfPublished string          `json:"year_of_published"  db:"year_of_published"`
	Rating          decimal.Decimal `json:"rating" db:"rating"`
}

type BookPhoto struct {
	PhotoName string
	Path      string
}

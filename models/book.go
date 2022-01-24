package library_app

import (
	"github.com/shopspring/decimal"
)

type Book struct {
	Id               int             `json:"id" db:"id"`
	BookTitle        string          `json:"book_title" db:"book_title"`
	BookTitleNative  string          `json:"book_title_native" db:"book_title_native"`
	BookPrice        decimal.Decimal `json:"book_price" db:"book_price"`
	InventoryCount   int             `json:"inventory_count" db:"inventory_count"`
	TotalCount       int             `json:"books_in_lib" db:"books_in_lib"`
	OneDayPrice      decimal.Decimal `json:"one_day_price" db:"one_day_price"`
	YearOfPublished  string          `json:"year_of_published"  db:"year_of_published"`
	RegistrationDate string          `json:"registration_date"  db:"registration_date"`
	NumberOfPages    int             `json:"number_of_pages" db:"number_of_pages"`
	BookState        bool            `json:"book_state" db:"book_state"`
	HideBook         bool            `json:"hide_book" db:"hide_book"`
	Rating           float64         `json:"rating" db:"rating"`
	Authors          []*Authors      `json:"Authors"`
	BooksPhoto       []*BooksPhoto   `json:"BooksPhoto"`
	Genres           []*Genres       `json:"Genres"`

}

type BooksPhoto struct {
	BookPhoto string `json:"book_photo" db:"book_photo"`
}

type Authors struct {
	Id           int             `json:"authors_id" db:"id"`
	FirstName    string          `json:"authors_firstname" db:"author_firstname"`
	LastName     string          `json:"authors_lastname" db:"author_lastname"`

}

type AuthorsPhoto struct {
	AuthorPhoto string `json:"authors_photo" db:"authors_photo"`
}

type Genres struct {
	Genre string `json:"genre" db:"genre"`
}

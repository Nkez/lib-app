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
	Rating           int             `json:"rating" db:"rating"`
	BookFoto         string          `json:"book_foto" db:"book_foto"`
	Authors          []struct {
		Id        int    `json:"authors_id" db:"id"`
		FirstName string `json:"authors_firstname" db:"authors_firstname"`
		LastName  string `json:"authors_lastname" db:"authors_lastname"`
		AutorFoto string `json:"authors_foto" db:"authors_foto"`
	} `json:"authors" db:"authors"`

	Genres struct {
		ActionAndAdventure string `json:"action"`
		//Classic             string "Classics"
		//DetectiveAndMystery string "Detective and Mystery"
		//Fantasy             string "Fantasy"
		//HistoricalFiction   string "Historical Fiction"
		//Horror              string "Horror"
	} `json:"genres"`
}

type ReturnBook struct {
	BookTitle       string `json:"book_title" db:"book_title"`
	Genre           string `json:"genre" db:"genre"`
	YearOfPublished string `json:"year_of_published"  db:"year_of_published"`
	InventoryCount  int    `json:"inventory_count" db:"inventory_count"`
	TotalCount      int    `json:"books_in_lib" db:"books_in_lib"`
}

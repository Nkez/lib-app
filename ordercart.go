package library_app

import (
	"time"
)

type OrderCart struct {
	LastName     string    `json:"last_name" db:"last_name"`
	FirstName    string    `json:"first_name" db:"first_name"`
	EmailAddress string    `json:"email_address" db:"email_address"`
	Price        float64   `json:"price" db:"price"`
	DateToReturn time.Time `json:"date_to_return " db:"date_to_return"`
	Book1        string    `json:"book1" db:"book1"`
	Book2        string    `json:"book2" db:"book2"`
	Book3        string    `json:"book3" db:"book3"`
	Book4        string    `json:"book4" db:"book4"`
	Book5        string    `json:"book5" db:"book5"`
}

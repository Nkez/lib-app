package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Order struct {
	Id           int             `json:"id"`
	FirstName    string          `json:"first_name" valid:"required"`
	LastName     string          `json:"last_name" valid:"required"`
	EmailAddress string          `json:"email_address" db:"email_address" valid:"email"`
	OrderDate    string          `json:"order_date"`
	DateToReturn string          `json:"date_to_return"`
	Price        decimal.Decimal `json:"price"`
	IsDebtor     bool            `json:"-"`
	Books        []string        `json:"books"`
}

type Return struct {
	FirstName    string
	LastName     string
	EmailAddress string
	OrderDate    string
	DateToReturn string
	Price        decimal.Decimal
	Defect       string
	DefectFoto   string
	IsDefect     bool
	Rating       decimal.Decimal
	IsReturn     bool
}

type OrderInput struct {
	EmailAddress string `json:"email_address" db:"email_address"valid:"email"`
	Books        []struct {
		Book string `json:"book" db:"book"`
	} `json:"books" db:"books"`
}

type ReturnInput struct {
	IdUser     int        `json:"id" db:"id"`
	OrderDay   string     `json:"order_date"`
	ReturnDay  string     `json:"return_day"`
	ReturnCart []ReturnST `json:"return_cart"`
}
type ReturnST struct {
	IdBook   int             `json:"book" db:"book"`
	IdCopies int             `json:"id_copies" db:"id_copies"`
	Rating   decimal.Decimal `json:"rating"`
	Defect   string          `json:"defect" db:"defect"`
}
type DefectFotos struct {
	IdDefectFoto int `json:"defect_foto"`
}

type InfoOrdDept struct {
	Id           int             `json:"id" db:"id"`
	FirstName    string          `json:"first_name" db:"first_name"`
	LastName     string          `json:"last_name" db:"last_name"`
	EmailAddress string          `json:"email_address" db:"email_address"`
	Price        decimal.Decimal `json:"price" db:"price"`
	OrderDate    string          `json:"order_date" db:"order_date"`
	DateToReturn string          `json:"date_to_return" db:"date_to_return"`
	Books        string
}

type DeborBooks struct {
	Book string
}

type DefetBooks struct {
	Id         int
	Defect     string
	DefectFoto string
}

type DbrInfo struct {
	Id           int       `json:"users_id" db:"users_id"`
	Price        float64   `json:"price" db:"price"`
	OrderDate    time.Time `json:"order_date" db:"order_date"`
	DateToReturn time.Time `json:"date_to_return" db:"date_to_return"`
}

type ReturnOrder struct {
	Id           int             `json:"id" db:"id"`
	FirstName    string          `json:"first_name" db:"first_name"`
	LastName     string          `json:"last_name" db:"last_name"`
	EmailAddress string          `json:"email_address" db:"email_address"`
	Price        decimal.Decimal `json:"price" db:"price"`
	OrderDate    string          `json:"order_date" db:"order_date"`
	DateToReturn string          `json:"date_to_return" db:"date_to_return"`
	Books        []string        `json:"books" db:"books"`
	IdCopies     []int           `json:"id_copies" db:"id_copies"`
}

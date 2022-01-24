package library_app

import "github.com/shopspring/decimal"

type UserBook struct {
	OrderDate string
	DateToReturn string
	Price decimal.Decimal
	Defect string
	IsDefect string
	Rating float64
}

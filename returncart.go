package library_app

import "github.com/shopspring/decimal"

//CREATE TABLE return_cart
//(
//date_to_return timestamp  not null default CURRENT_TIMESTAMP,
//price DECIMAL(6,2) NOT NULL,
//rating INT,
//defect_foto VARCHAR(255),
//defect VARCHAR(255),
//is_book_defect BOOLEAN
//);

type ReturnCart struct {
	ReturnDate        string
	Price             decimal.Decimal
	Rating            int
	DefectFoto        string
	DefectDescription string
	IsDefect          bool
}

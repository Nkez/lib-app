package repository

import (
	"database/sql"
	"errors"
	"fmt"
	library_app "github.com/Nkez/lib-app.git"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type CartPostgres struct {
	db *sqlx.DB
}

func NewCartPostgres(db *sqlx.DB) *CartPostgres {
	return &CartPostgres{db: db}
}

func (r *CartPostgres) CreateCart(user library_app.OrderCart) (library_app.OrderCart, error) {

	query := fmt.Sprintf(`INSERT INTO %s(last_name,first_name,email_address,
								book1,book2,book3,book4,book5,price,date_to_return)
								VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING email_address`, cartTable)
	if err := r.db.QueryRow(query, user.LastName, user.FirstName, user.EmailAddress,
		user.Book1, user.Book2, user.Book3, user.Book4, user.Book5,
		user.Price, user.DateToReturn); err != nil {
		return user, nil
	}
	return user, nil
}

func (r *CartPostgres) FindRegisterUser(email string) (library_app.User, error) {
	user := library_app.User{}
	logrus.Info("Find Register User DB")
	query := fmt.Sprintf(`SELECT id,email_address,last_name,first_name FROM %s WHERE  email_address = $1`, userTable)
	if err := r.db.QueryRow(query,
		email,
	).Scan(
		&user.Id,
		&user.EmailAddress,
		&user.LastName,
		&user.FirstName,
	); err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("Create new user")
		}
		return user, err
	}
	return user, nil
}

func (r *CartPostgres) CheckBooks(book string) (checkBook string, err error) {
	logrus.Info("Check book DB")
	query := fmt.Sprintf(`SELECT book_title FROM %s WHERE book_title = $1 AND inventory_count> 0 AND hide_book=FALSE`, bookTable)
	if err := r.db.QueryRow(query, book).Scan(&checkBook); err != nil {
		return "unknown", fmt.Errorf("didnt find book: %q", book)
	}

	r.RemoveInventoryCount(book)
	return checkBook, nil
}

func (r *CartPostgres) GetPrice(book string) (priceBook float64, err error) {
	logrus.Info("Price User DB")
	query := fmt.Sprintf(`SELECT book_price FROM %s WHERE book_title = $1`, bookTable)
	if err := r.db.QueryRow(query, book).Scan(&priceBook); err != nil {
		logrus.Info("FORGET ABOUT PRICE")
		return 0, errors.New("chek price in GetPrice DB")
	}
	return priceBook, nil
}

func (r *CartPostgres) GetUser(email string) (library_app.OrderCart, error) {
	logrus.Info("Get User DB")
	cart := library_app.OrderCart{}

	query := fmt.Sprintf(`SELECT first_name,email_address,last_name FROM %s WHERE  email_address = $1`, userTable)
	if row := r.db.QueryRow(query, email).Scan(&cart.FirstName, &cart.EmailAddress, &cart.LastName); row != nil {
		return cart, errors.New("USER NOT FOUND")
	}

	return cart, nil
}

func (r *CartPostgres) RemoveInventoryCount(book string) {
	logrus.Info("UPDATE inventory_count inventory_count")
	queryUpdate := fmt.Sprintf(`UPDATE %s SET inventory_count=inventory_count-1 WHERE book_title = $1`, bookTable)
	r.db.QueryRow(queryUpdate, book)
}

type CheckBooks struct {
	Book1 string
	Book2 string
	Book3 string
	Book4 string
	Book5 string
}

func (r *CartPostgres) CheckOrderBook(email string) (orderBook []string, err error) {
	logrus.Info("Check order book DB")
	books := CheckBooks{}
	query := fmt.Sprintf(`SELECT book1, book2,book3, book4, book5 FROM %s WHERE email_address=$1`, cartTable)
	err = r.db.QueryRow(query, email).Scan(&books.Book1, &books.Book2, &books.Book3, &books.Book4, &books.Book5)
	if err != nil {
		return nil, fmt.Errorf("RETURN BOOOOK!")
	}
	orderBook = append(orderBook, books.Book1)
	orderBook = append(orderBook, books.Book2)
	orderBook = append(orderBook, books.Book3)
	orderBook = append(orderBook, books.Book4)
	orderBook = append(orderBook, books.Book5)
	return orderBook, nil
}

func (r *CartPostgres) UpdatePrice() {
	logrus.Info("Update Price")
	query := fmt.Sprintf(`UPDATE %s SET price=price*2 WHERE date_to_return < NOW()`, cartTable)
	_, err := r.db.Query(query)
	if err != nil {
		fmt.Errorf("problem in updatePrice")
	}
}

func (r *CartPostgres) GetEmailToSend() ([]library_app.OrderCart, error) {
	logrus.Info("Get email to send")

	var carts []library_app.OrderCart
	query := fmt.Sprintf(`SELECT * FROM %s WHERE date_to_return < NOW()`, cartTable)
	if err := r.db.Select(&carts, query); err != nil {
		return carts, err
	}
	return carts, nil
}

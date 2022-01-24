package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Nkez/library-app.git/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) FindUser(email string) (models.User, error) {
	user := models.User{}
	logrus.Info("Find Register User DB")
	query := fmt.Sprintf(`SELECT id,email_address,last_name,first_name FROM %s
								WHERE  email_address = $1`, userTable)
	if err := r.db.QueryRow(query,
		email,
	).Scan(
		&user.Id,
		&user.EmailAddress,
		&user.LastName,
		&user.FirstName,
	); err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("create new user")
		}
		return user, err
	}
	return user, nil
}

func (r *OrderPostgres) GetBooksId(books []string) (idBooks []int, err error) {
	logrus.Info("Check book DB and Price")
	var id int
	query := fmt.Sprintf(`SELECT id FROM %s 
								WHERE book_title = $1 AND inventory_count> 0`, bookTable)
	for i := 0; i < len(books); i++ {
		if err := r.db.QueryRow(query, books[i]).Scan(&id); err != nil {
			return nil, errors.New(fmt.Sprintf("didnt find book %s", books[i]))

		}
		idBooks = append(idBooks, id)
	}
	return idBooks, nil
}

func (r *OrderPostgres) GetPrice(id []int) (prices []float64, err error) {
	logrus.Info("Check book DB and Price")
	var price float64
	query := fmt.Sprintf(`SELECT one_day_price FROM %s 
								WHERE id = $1 AND inventory_count> 0`, bookTable)
	for i := 0; i < len(id); i++ {
		if err := r.db.QueryRow(query, id[i]).Scan(&price); err != nil {
			return nil, errors.New("error in data base scan ")

		}
		prices = append(prices, price)
	}
	return prices, nil
}

func (r *OrderPostgres) CheckIsReturn(id int) (books []string, err error) {
	logrus.Info("Check is return")
	query := fmt.Sprintf(`SELECT  b.book_title 
								FROM %s bu
								INNER JOIN %s as u
									ON u.id = bu.users_id
								INNER JOIN %s b 
									ON b.id = bu.books_id 
								WHERE u.id = $1 and bu.is_return = false  `, bookUserTable, userTable, bookTable)
	if err := r.db.Select(&books, query, id); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *OrderPostgres) JoinBookUser(idUser int, idBook []int, price float64, orderDate, returnDate time.Time) error {
	logrus.Info("join  user book DB")
	isReturn := false
	query := fmt.Sprintf(`INSERT INTO %s(order_date,date_to_return ,price,is_return ,books_id ,users_id )
								 VALUES($1,$2,$3,$4,$5,$6)`, bookUserTable)
	for i := 0; i < len(idBook); i++ {
		r.db.QueryRow(query, orderDate, returnDate, price, isReturn, idBook[i], idUser)

	}
	//queryUp := fmt.Sprintf(`update books_users set is_return= fasle where users_id = $1`)
	//r.db.QueryRow(queryUp, idUser)
	return nil
}

func (r *OrderPostgres) MinusInventoryCount(idBooks []int) error {
	logrus.Info("Minus inventory count BD")
	query := fmt.Sprintf(`UPDATE %s SET inventory_count=inventory_count-1 WHERE id = $1`, bookTable)
	fmt.Println("invenroty count")
	fmt.Println(idBooks)
	for i := 0; i < len(idBooks); i++ {
		r.db.QueryRow(query, idBooks[i])
	}
	return nil
}

func (r *OrderPostgres) ReturnOrder(id int) ([]string, models.Order, error) {
	logrus.Info("Return order cart DB")
	var order models.Order
	query := fmt.Sprintf(`select  u.id,u.first_name , u.last_name , u.email_address, bu.order_date, 
								bu.date_to_return , bu.price 
								from books_users bu
								inner join users as u
									on u.id = bu.users_id 
								where u.id = $1 and bu.is_return = false`)
	if err := r.db.QueryRow(query,
		id,
	).Scan(
		&order.FirstName,
		&order.LastName,
		&order.EmailAddress,
		&order.OrderDate,
		&order.DateToReturn,
		&order.Price,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, order, errors.New("create new user")
		}
		return nil, order, err
	}

	logrus.Info("get books book-user DB")
	var books []string
	queryBook := fmt.Sprintf(` select  b.book_title 
									from books_users bu
									inner join users as u
										on u.id = bu.users_id 
									inner join books as b
										on b.id = bu.books_id 
									where u.id = $1 and bu.is_return = false`)
	if err := r.db.Select(&books, queryBook, id); err != nil {
		return nil, order, err
	}
	fmt.Println(books)
	return books, order, nil
}

func (r *OrderPostgres) GetAllOrder(page, limit string) ([]models.InfoOrdDept, error) {
	logrus.Info("Return order cart DB")
	var order []models.InfoOrdDept
	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)
	pageForSql := (p - 1) * 5
	query := fmt.Sprintf(` select distinct  u.id ,u.first_name ,u.last_name ,u.email_address,
									bu.price, bu.order_date ,bu.date_to_return 
									from books_users bu
									left join users as u
										on u.id = bu.users_id 
									left join books as b
										on b.id = bu.books_id
										where bu.is_return = false LIMIT %d OFFSET %d`, l, pageForSql)
	if err := r.db.Select(&order, query); err != nil {
		return nil, err
	}

	queryBook := fmt.Sprintf(`select  b.book_title 
									from books_users bu
									inner join users as u
										on u.id = bu.users_id 
									inner join books as b
										on b.id = bu.books_id			
										where u.id  = $1 and bu.is_return = false`)
	for i, item := range order {
		var books []string
		r.db.Select(&books, queryBook, item.Id)
		order[i].Books = strings.Join(books, ",")
	}

	return order, nil
}

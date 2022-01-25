package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Nkez/library-app.git/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) FindUser(email string) (models.User, error) {
	logrus.Info(fmt.Sprintf("find user, email: %s", email))
	user := models.User{}
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
	logrus.Info(fmt.Sprintf("get books price, %s", books))
	query := fmt.Sprintf(`SELECT id FROM %s 
								WHERE book_title = any($1) AND inventory_count> 0`, bookTable)
	if err = r.db.Select(&idBooks, query, pq.Array(books)); err != nil {
		return nil, err
	}

	return idBooks, nil
}

func (r *OrderPostgres) GetPrice(id []int) (prices []float64, err error) {
	logrus.Info(fmt.Sprintf("check book price, id: %v", id))
	query := fmt.Sprintf(`SELECT one_day_price FROM %s 
								WHERE id = any($1) AND inventory_count> 0`, bookTable)
	if err = r.db.Select(&prices, query, pq.Array(id)); err != nil {
		return nil, err
	}
	return prices, nil
}

func (r *OrderPostgres) CheckIsReturn(id int) (books []string, err error) {
	logrus.Info(fmt.Sprintf("check is return user, id :%v", id))
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

func (r *OrderPostgres) JoinBookCopies(idBook []int) (idCopies []int, err error) {
	logrus.Info(fmt.Sprintf("join books id %v,book copies", idBook))

	query := fmt.Sprintf(`insert into book_copies (books_id) values ($1) returning id`)
	var id int
	for _, v := range idBook {
		row := r.db.QueryRow(query, v)
		if err = row.Scan(&id); err != nil {
			return nil, err
		}
		idCopies = append(idCopies, id)
	}
	return idCopies, nil
}

func (r *OrderPostgres) JoinBookUser(idUser int, idBooks []int, price float64, orderDate, returnDate time.Time) error {
	logrus.Info(fmt.Sprintf("join books : %v, users : %v", idBooks, idUser))

	isReturn := false
	query := fmt.Sprintf(`INSERT INTO %s(order_date,date_to_return ,price,is_return ,books_id ,users_id )
								 VALUES($1,$2,$3,$4,$5,$6)`, bookUserTable)
	for i := 0; i < len(idBooks); i++ {
		r.db.QueryRow(query, orderDate, returnDate, price, isReturn, idBooks[i], idUser)
	}
	return nil
}

func (r *OrderPostgres) MinusInventoryCount(idBooks []int) error {
	logrus.Info("Minus inventory count BD")
	query := fmt.Sprintf(`UPDATE %s SET inventory_count=inventory_count-1 WHERE id = $1`, bookTable)
	for i := 0; i < len(idBooks); i++ {
		r.db.QueryRow(query, idBooks[i])
	}
	return nil
}

func (r *OrderPostgres) ReturnOrder(id int) (order models.ReturnOrder, err error) {
	logrus.Info(fmt.Sprintf("return order, id: %v", id))

	query := fmt.Sprintf(` select u.id,u.first_name, u.last_name ,u.email_address , bu.order_date, bu.date_to_return ,bu.price,
								b.book_title , bc.id from books b
								join book_copies bc on bc.id  = b.id  
								join books_users bu on bu.books_id  = b.id  
								join users u on u.id = bu.users_id 
								where u.id = $1 and bu.is_return = false `)
	row, err := r.db.Query(query, id)

	var books []string
	var idCopies []int
	defer row.Close()
	for row.Next() {
		var b string
		var ids int
		if err := row.Scan(&order.Id, &order.FirstName, &order.LastName, &order.EmailAddress, &order.OrderDate, &order.DateToReturn, &order.Price, &b, &ids); err != nil {
			return order, nil
		}

		books = append(books, b)
		idCopies = append(idCopies, ids)
	}

	order.Books = books
	order.IdCopies = idCopies

	return order, err
}

func (r *OrderPostgres) GetAllOrder(page, limit string) ([]models.ReturnOrder, error) {
	logrus.Info("Return order cart DB")

	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)
	pageForSql := (p - 1) * 5
	query := fmt.Sprintf(` select u.id,u.first_name, u.last_name ,u.email_address , bu.order_date, bu.date_to_return ,bu.price
								from books b
								join book_copies bc on bc.books_id  = b.id  
								join books_users bu on bu.books_id  = b.id  
								join users u on u.id = bu.users_id 
								where bu.is_return = false LIMIT %d OFFSET %d`, l, pageForSql)
	row, _ := r.db.Query(query)

	defer row.Close()
	var orders []models.ReturnOrder
	for row.Next() {
		order := models.ReturnOrder{}

		if err := row.Scan(&order.Id, &order.FirstName, &order.LastName, &order.EmailAddress, &order.OrderDate, &order.DateToReturn, &order.Price); err != nil {
			return orders, nil
		}

		order.IdCopies = r.GetCopies(order.Id)
		order.Books = r.GetBooks(order.Id)
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderPostgres) GetCopies(id int) (copies []int) {
	logrus.Info(fmt.Sprintf("get copies, id %v", id))
	quert := fmt.Sprintf(`select 
								 bc.id from books b
								join book_copies bc on bc.books_id  = b.id  
								join books_users bu on bu.books_id  = b.id  
								join users u on u.id = bu.users_id 
								where u.id = $1 and bu.is_return = false`)
	r.db.Select(&copies, quert, id)
	return copies
}

func (r *OrderPostgres) GetBooks(id int) (books []string) {
	logrus.Info(fmt.Sprintf("get books, id %v", id))
	quert := fmt.Sprintf(`select 
								 b.book_title from books b
								join book_copies bc on bc.books_id  = b.id  
								join books_users bu on bu.books_id  = b.id  
								join users u on u.id = bu.users_id 
								where u.id = $1 and bu.is_return = false`)
	r.db.Select(&books, quert, id)
	return books
}

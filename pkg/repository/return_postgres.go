package repository

import (
	"fmt"
	"github.com/Nkez/library-app.git/models"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ReturnPostgres struct {
	db *sqlx.DB
}
type DeborsPostgres struct {
	db *sqlx.DB
}

func NewReturnPostgres(db *sqlx.DB) *ReturnPostgres {
	return &ReturnPostgres{db: db}
}
func NewDeborsPostgres(db *sqlx.DB) *DeborsPostgres {
	return &DeborsPostgres{db: db}
}

func (r *DeborsPostgres) GetAllDebors() ([]models.InfoOrdDept, error) {
	logrus.Info("get all debors from db")
	var debors []models.InfoOrdDept
	boolF := false
	query := fmt.Sprintf(` select distinct  u.id ,u.first_name ,u.last_name ,u.email_address,
									bu.price, bu.order_date ,bu.date_to_return 
									from books_users bu
									left join users as u
										on u.id = bu.users_id 
									left join books as b
										on b.id = bu.books_id
										where is_return= $1 and bu.date_to_return < NOW() LIMIT 20`)
	if err := r.db.Select(&debors, query, boolF); err != nil {
		return nil, err
	}

	queryBook := fmt.Sprintf(`select  b.book_title 
									from books_users bu
									inner join users as u
										on u.id = bu.users_id 
									inner join books as b
										on b.id = bu.books_id			
										where u.id  = $1`)
	for i, item := range debors {
		var books []string
		r.db.Select(&books, queryBook, item.Id)
		debors[i].Books = strings.Join(books, ",")
	}
	return debors, nil
}

func (r *ReturnPostgres) CheckReturningBook(id int) ([]int, error) {
	logrus.Info("check returning book DB")
	var idCheck []int
	query := fmt.Sprintf(`select b.id 
									from books_users bu
									inner join users as u
										on u.id = bu.users_id 
									inner join books as b
										on b.id = bu.books_id 
									where u.id = $1 and bu.is_return = false `)
	r.db.Select(&idCheck, query, id)
	return idCheck, nil
}

func (r *ReturnPostgres) ReturnBook(idUser int, idBook []int, returnDate time.Time) (models.DbrInfo, error) {
	logrus.Info("return book db")
	var dbrInfo models.DbrInfo
	rtBool := true

	query := fmt.Sprintf(`update  books_users  SET  is_return = $2, order_date = $3  WHERE users_id = $1 RETURNING users_id, order_date,
							date_to_return, price`)
	if err := r.db.QueryRow(query, idUser, rtBool, returnDate).Scan(&dbrInfo.Id, &dbrInfo.OrderDate, &dbrInfo.DateToReturn, &dbrInfo.Price); err != nil {
		return dbrInfo, err
	}

	return dbrInfo, nil
}

func (r *ReturnPostgres) ReturnBooksInLibUpdateRating(booksId []int, rating []decimal.Decimal) error {
	logrus.Info("UPDATE ok books DB")
	queryUpdate := fmt.Sprintf(`UPDATE %s SET inventory_count=inventory_count+1, rating = $1 WHERE id = $2 `, bookTable)
	for i := 0; i < len(booksId); i++ {
		r.db.QueryRow(queryUpdate, rating[i], booksId[i])
	}
	return nil
}

func (r *ReturnPostgres) GetRating(idBooks []int) ([]float64, error) {
	logrus.Info("get book rating")
	var rating float64
	var allRatings []float64
	query := fmt.Sprintf("select rating from %s where id = $1", bookTable)
	for i := 0; i < len(idBooks); i++ {
		row := r.db.QueryRow(query, idBooks[i])
		if err := row.Scan(&rating); err != nil {
			allRatings = append(allRatings, 0.0)
		}
		allRatings = append(allRatings, rating)
	}
	return allRatings, nil
}

func (r *ReturnPostgres) UpdatePrice(id int, newPrice float64) {
	query := fmt.Sprintf(`UPDATE books_users SET price = $2 WHERE users_id = $1`)
	r.db.QueryRow(query, id, newPrice)
}

//func(r *ReturnPostgres) JoinDefectBook(idBook int, paths []string) error{
//	logrus.Info("join book defect photo")
//	query := fmt.Sprintf(`insert into books_books_photo values($1,$2)`)
//	for i := 1; i < len(paths)+1; i++ {
//		err := r.db.QueryRow(query, idBook,	i)
//		if err != nil {
//			return fmt.Errorf("error in join books")
//		}
//	}
//	return nil
//}

/*logrus.Info("get rating  DB")
var uCount int
var rating float64
queryUpdate := fmt.Sprintf(` select  count( u.last_name), b.rating
										from books_users bu
										inner join users as u
											on u.id = bu.users_id
										inner join books as b
										on b.id = bu.books_id
										where  b.id= $1
										group by
	 									b.rating `)

if err := r.db.QueryRow(queryUpdate, book).Scan(&uCount, &rating); err != nil {
return 0, 0.0, errors.New("error get rating DB")
}

return uCount, rating, nil*/

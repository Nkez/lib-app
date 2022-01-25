package repository

import (
	"fmt"
	"github.com/Nkez/library-app.git/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type UserList struct {
	id             int
	lastName       string
	firstName      string
	middleName     string
	passportNumber string
	birthday       time.Time
	emailAddress   string
	address        string
}

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) CreateUser(user models.User) (int, error) {
	logrus.Info("create user")
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (last_name, first_name,middle_name,passport_number,birthday,email_address,address)
								VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`, userTable)
	row := r.db.QueryRow(query, user.LastName, user.FirstName, user.MiddleName, user.PassportNumber, user.Birthday, user.EmailAddress, user.Address)
	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func (r *UserPostgres) GetAllUsers(page, limit string) ([]models.User, error) {
	logrus.Info("find all users")
	var users []models.User
	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)
	pageForSql := (p - 1) * 5
	query := fmt.Sprintf(`SELECT last_name,first_name, birthday, email_address FROM users  LIMIT %d OFFSET %d`, l, pageForSql)
	if err := r.db.Select(&users, query); err != nil {
		return users, err
	}
	return users, nil
}

func (r *UserPostgres) GetByName(name string) []models.UserName {
	logrus.Info(fmt.Sprintf("find user, name%s", name))
	var users []models.UserName
	query := fmt.Sprintf("SELECT last_name,first_name,email_address FROM %s where first_name = $1", userTable)
	r.db.Select(&users, query, name)
	return users
}

func (r *UserPostgres) FindByWord(name string) []models.UserName {
	logrus.Info(fmt.Sprintf("find user, word%s", name))
	var s []models.UserName
	query := fmt.Sprintf(`SELECT first_name,last_name,email_address FROM users WHERE last_name LIKE $1 or first_name LIKE $1`)
	r.db.Select(&s, query, name+"%")
	return s
}

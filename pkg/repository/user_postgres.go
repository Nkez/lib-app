package repository

import (
	"fmt"
	library_app "github.com/Nkez/lib-app.git"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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

func (r *UserPostgres) CreateUser(user library_app.User) (int, error) {
	logrus.Info("Create User DB")
	var id int
	query := fmt.Sprintf("INSERT INTO %s (last_name, first_name,middle_name,passport_number,birthday,email_address,address) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id", userTable)
	row := r.db.QueryRow(query, user.LastName, user.FirstName, user.MiddleName, user.PassportNumber, user.Birthday, user.EmailAddress, user.Address)
	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func (r *UserPostgres) GetAllUsers() ([]library_app.User, error) {
	logrus.Info("Find All Users DB")
	var users []library_app.User

	query := fmt.Sprintf("SELECT last_name,first_name, birthday, email_address FROM %s", userTable)
	if err := r.db.Select(&users, query); err != nil {
		return users, err
	}
	return users, nil

}

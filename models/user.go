package models

type User struct {
	Id             int    `json:"-" db:"id"`
	LastName       string `json:"last_name" db:"last_name"  valid:"required"`
	FirstName      string `json:"first_name" db:"first_name"valid:"required"`
	MiddleName     string `json:"middle_name" db:"middle_name"valid:"required"`
	PassportNumber string `json:"passport_number" db:"passport_number"`
	Birthday       string `json:"birthday" db:"birthday"`
	EmailAddress   string `json:"email_address" db:"email_address"valid:"required,email"`
	Address        string `json:"address" db:"address"`
}

type UserName struct {
	LastName  string `json:"last_name" db:"last_name"  valid:"required"`
	FirstName string `json:"first_name" db:"first_name"valid:"required"`
	EmailAddr string `json:"email_address" db:"email_address"valid:"required,email"`
}

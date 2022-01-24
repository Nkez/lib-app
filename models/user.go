package library_app

type User struct {
	Id             int    `json:"-" db:"id"`
	LastName       string `json:"last_name" db:"last_name"`
	FirstName      string `json:"first_name" db:"first_name"`
	MiddleName     string `json:"middle_name" db:"middle_name"`
	PassportNumber string `json:"passport_number" db:"passport_number"`
	Birthday       string `json:"birthday" db:"birthday"`
	EmailAddress   string `json:"email_address" db:"email_address"`
	Address        string `json:"address" db:"address"`
}

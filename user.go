package library_app

type User struct {
	Id             int    `json:"Id"`
	LastName       string `json:"Last_name"`
	FirstName      string `json:"First_name"`
	MiddleName     string `json:"Middle_name"`
	PasswordNumber string `json:"Password_number"`
	EmailAddress   string `json:"Email_address"`
	Birthday       int    `json:"Birthday"`
	Address        string `json:"Address"`
}

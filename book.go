package library_app

type Author struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

type Book struct {
	BookTitle        string   `json:"BookTitle"`
	BookTitleNative  string   `json:"BookTitleNative"`
	Genre            string   `json:"Genre"`
	BookPrice        float32  `json:"BookPrice"`
	InventoryCount   int      `json:"Inventory_count"`
	Authors          []Author `json:"Authors"`
	NumberOfPages    int      `json:"NumberOfPages"`
	RegistrationDate string   `json:"RegistrationDate"`
}

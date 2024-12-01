package models

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Order struct {
	UserID     string   `json:"user_id"`
	ProductIDs []string `json:"product_ids"`
}

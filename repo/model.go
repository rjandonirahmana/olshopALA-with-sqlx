package repo

import "database/sql"

type Customer struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    int64  `json:"phone"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type Product struct {
	Name  string `json:"name"`
	ID    int    `json:"product_id"`
	Price int    `json:"price"`
}

type Transactions struct {
	ID_trans   int `json:"id" db:"id"`
	ID_product int `json:"product_id" db:"product_id"`
	CustomerID int `json:"customer_id" db:"customer_id"`
	Quantity   int `json:"quantity" db:"quantity"`
	Price      int `json:"price" db:"price"`
}

type CustomerDB struct {
	ID       sql.NullInt32  `db:"id"`
	Name     sql.NullString `db:"name"`
	Phone    sql.NullInt64  `db:"phone"`
	Email    sql.NullString `db:"email"`
	Password sql.NullString `db:"password"`
	Salt     sql.NullString `db:"salt"`
}

type ProductImage struct {
	ID         int    `json:"id"`
	ProductID  int    `json:"product_id"`
	IsPrimary  int    `json:"is_primary"`
	ImagesName string `json:"images_name"`
}

// type User struct {
// 	Id       int
// 	Name     string
// 	Gender   string
// 	Password string
// 	Email    string
// 	Token    string
// }

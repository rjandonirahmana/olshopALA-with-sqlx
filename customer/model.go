package customer

import "database/sql"

type Customer struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    int64  `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Salt     string `json:"salt"`
	Avatar   string `json:"avatar"`
}

type CustomerDB struct {
	ID       sql.NullInt32  `db:"id"`
	Name     sql.NullString `db:"name"`
	Phone    sql.NullInt64  `db:"phone"`
	Email    sql.NullString `db:"email"`
	Password sql.NullString `db:"password"`
	Salt     sql.NullString `db:"salt"`
	Avatar   sql.NullString `db:"avatar"`
}

type InputCustomer struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type InputLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

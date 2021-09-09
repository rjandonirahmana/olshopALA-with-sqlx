package customer

import (
	"database/sql"
	"time"
)

type Customer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	Salt      string    `json:"salt"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CustomerDB struct {
	ID        sql.NullInt32  `db:"id"`
	Name      sql.NullString `db:"name"`
	Phone     sql.NullString `db:"phone"`
	Email     sql.NullString `db:"email"`
	Password  sql.NullString `db:"password"`
	Salt      sql.NullString `db:"salt"`
	Avatar    sql.NullString `db:"avatar"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

type InputCustomer struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type InputLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

package seller

import "time"

type Seller struct {
	ID        uint      `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Phone     *string   `db:"phone" json:"phone_number"`
	Password  string    `db:"password" json:"-"`
	Salt      string    `db:"salt" json:"-"`
	Avatar    *string   `db:"avatar" json:"avatar"`
	Address   *string   `db:"adress" json:"address"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type InputSeller struct {
	Name            string `json:"name" binding:"required" validate:"required"`
	Email           string `json:"email" binding:"required" validate:"email"`
	Password        string `json:"password" binding:"required" validate:"min=8,max=32,alphanum"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"eqfield=Password,required"`
}

type InputLoginSeller struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required"`
}

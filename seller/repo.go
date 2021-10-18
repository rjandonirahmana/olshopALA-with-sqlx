package seller

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

type RepoInt interface {
	CreateSeller(seller Seller) (Seller, error)
	ChangePassword(newPassword string, id int) error
	GetSellerByEmail(email string) (Seller, error)
	GetSellerByID(id int) (Seller, error)
}

func NewRepository(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateSeller(seller Seller) (Seller, error) {
	querry := `INSERT INTO seller (name, phone, email, password, salt, avatar, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(querry, seller.Name, seller.Phone, seller.Email, seller.Password, seller.Salt, seller.Avatar, seller.CreatedAt, seller.UpdatedAt)

	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

func (r *Repo) ChangePassword(newPassword string, id int) error {
	querry := `UPDATE seller SET password = ? WHERE id = ? `

	_, err := r.db.Exec(querry, newPassword, id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *Repo) GetSellerByEmail(email string) (Seller, error) {
	querry := `SELECT * FROM seller WHERE email = ?`

	var seller Seller

	err := r.db.Get(&seller, querry, email)
	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

func (r *Repo) GetSellerByID(id int) (Seller, error) {
	querry := `SELECT * FROM seller WHERE id = ?`

	var seller Seller

	err := r.db.Get(&seller, querry, id)
	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

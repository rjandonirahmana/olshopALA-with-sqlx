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
	GetLastID() (int, error)
	ChangePassword(newPassword string, id int) error
	GetSellerByEmail(email string) (Seller, error)
}

func (r *Repo) CreateSeller(seller Seller) (Seller, error) {
	querry := `INSERT INTO seller (id, name, phone, email, password, salt, avatar, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(querry, seller.ID, seller.Name, seller.Phone, seller.Email, seller.Password, seller.Salt, seller.Avatar, seller.CreatedAt, seller.UpdatedAt)

	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

func (r *Repo) GetLastID() (int, error) {
	querry := `SELECT id FROM seller WHERE id = (SELECT MAX(id) FROM seller)`

	var value int
	err := r.db.Get(&value, querry)
	if err != nil {
		return 0, err
	}
	return value, nil

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

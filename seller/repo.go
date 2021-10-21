package seller

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

type RepoInt interface {
	CreateSeller(seller Seller) (Seller, error)
	ChangePassword(newPassword string, id uint) error
	GetSellerByEmail(email string) (Seller, error)
	GetSellerByID(id uint) (Seller, error)
	IsEmailAvailable(email string) error
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreateSeller(seller Seller) (Seller, error) {
	querry := `INSERT INTO seller (name, email, password, salt, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id uint
	err := r.db.QueryRowx(querry, seller.Name, seller.Email, seller.Password, seller.Salt, seller.CreatedAt, seller.UpdatedAt).Scan(&id)

	if err != nil {
		return Seller{}, err
	}
	seller.ID = id

	return seller, nil
}

func (r *repository) ChangePassword(newPassword string, id uint) error {
	querry := `UPDATE seller SET password = $1 WHERE id = $2`

	_, err := r.db.Exec(querry, newPassword, id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repository) GetSellerByEmail(email string) (Seller, error) {
	querry := `SELECT * FROM seller WHERE email = $1`

	var seller Seller

	err := r.db.Get(&seller, querry, email)
	if err == sql.ErrNoRows {
		return Seller{}, errors.New("email not found")
	}
	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

func (r *repository) GetSellerByID(id uint) (Seller, error) {
	querry := `SELECT * FROM seller WHERE id = $1`

	var seller Seller

	err := r.db.Get(&seller, querry, id)
	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

func (r *repository) IsEmailAvailable(email string) error {
	var id uint
	querry := `SELECT id FROM seller WHERE email = $1`
	err := r.db.QueryRowx(querry, email).Scan(&id)
	if err == sql.ErrNoRows && id == 0 {
		return nil
	}
	if err != nil {
		return err
	}

	return errors.New("email has been used by another user")
}

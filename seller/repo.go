package seller

import (
	"database/sql"
	"errors"
	"fmt"
	"olshop/product"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

type RepoInt interface {
	CreateSeller(seller Seller) (Seller, error)
	ChangePassword(newPassword string, id int) error
	GetSellerByEmail(email string) (Seller, error)
	IsEmailAvailable(email string) (bool, error)
	InsertProduct(product product.Product) (int64, error)
	ChangeQuantityProduct(product_id int, quantity int) error
	GetSellerByID(id int) (Seller, error)
	SelectTypeProductID(type_name string) (int, error)
}

func NewRepoSeller(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateSeller(seller Seller) (Seller, error) {
	querry := `INSERT INTO sellers (id, name, phone, email, password, salt, avatar, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(querry, seller.ID, seller.Name, seller.Phone, seller.Email, seller.Password, seller.Salt, seller.Avatar, seller.CreatedAt, seller.UpdatedAt)

	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

func (r *Repo) GetSellerByID(id int) (Seller, error) {
	querry := `SELECT * FROM sellers WHERE id = ?`

	var seller Seller
	err := r.db.Get(&seller, querry, id)
	if err != nil {
		return Seller{}, err
	}
	return seller, nil
}

func (r *Repo) ChangePassword(newPassword string, id int) error {
	querry := `UPDATE sellers SET password = ?, updated_at WHERE id = ? `

	_, err := r.db.Exec(querry, newPassword, time.Now(), id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *Repo) GetSellerByEmail(email string) (Seller, error) {
	querry := `SELECT * FROM sellers WHERE email = ?`

	var seller Seller

	err := r.db.Get(&seller, querry, email)
	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

func (r *Repo) IsEmailAvailable(email string) (bool, error) {
	querry := `SELECT * FROM sellers WHERE email = ?`

	var seller Seller

	err := r.db.Get(&seller, querry, email)
	if err != sql.ErrNoRows {
		return false, err
	}
	if seller.Email == email {
		return false, errors.New("email has been used")
	}
	return true, nil
}

func (r *Repo) InsertProduct(product product.Product) (int64, error) {
	querry := `INSERT INTO products(id, name, price, category_id, description, quantity, seller_id) VALUES (?,?,?,?,?,?,?)`

	sq, err := r.db.Exec(querry, product.ID, product.Name, product.Price, product.Category_id, product.Description, product.Quantity, product.SellerID)

	if err != nil {
		return 0, err
	}
	lastid, err := sq.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastid, nil
}

func (r *Repo) SelectTypeProductID(type_name string) (int, error) {
	querry := `SElECT id FROM product_category where name = ?`

	var id int
	err := r.db.Get(&id, querry, type_name)
	if err != nil {
		return 0, fmt.Errorf("product_id not found or err : %v", err)
	}

	return id, nil

}

func (r *Repo) ChangeQuantityProduct(product_id int, quantity int) error {
	querry := `UPDATE products SET quantity = ? WHERE id = ?`

	_, err := r.db.Exec(querry, quantity, product_id)
	if err != nil {
		return err
	}
	return nil

}

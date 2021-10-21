package customer

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

type Repository interface {
	RegisterUser(customer Customer) (Customer, error)
	UpdateCustomerPhone(id uint, number string) error
	GetCustomerByID(id uint) (Customer, error)
	ChangePassword(newPassword string, id uint) error
	GetCustomerByEmail(email string) (Customer, error)
	ChangeAvatar(avatarFile string, id uint) error
	DeleteCustomer(id uint) error
	IsEmailAvailable(email string) error
}

func NewRepo(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) RegisterUser(customer Customer) (Customer, error) {

	querry := `INSERT INTO customers (name, phone, email, password, salt, avatar, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	var id uint
	err := r.db.QueryRowx(querry, customer.Name, customer.Phone, customer.Email, customer.Password, customer.Salt, customer.Avatar, customer.CreatedAt, customer.UpdatedAt).Scan(&id)

	if err != nil {
		return Customer{}, err
	}

	customer.ID = id
	return customer, nil
}

func (r *repository) UpdateCustomerPhone(id uint, number string) error {

	querry := `
	UPDATE 
		customers 
	SET 
		phone = $1, 
		updated_at = $2
	WHERE 
		id = $3
	`

	_, err := r.db.Exec(querry, number, time.Now(), id)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetCustomerByID(id uint) (Customer, error) {
	querry := `SELECT * FROM customers WHERE id = $1`

	var customer Customer

	err := r.db.Get(&customer, querry, id)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (r *repository) ChangeAvatar(avatarFile string, id uint) error {
	querry := `UPDATE customers SET avatar = $1, updated_at = $2 WHERE id = $3`

	_, err := r.db.Exec(querry, avatarFile, time.Now(), id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repository) ChangePassword(newPassword string, id uint) error {
	querry := `UPDATE customers SET password = $1, updated_at = $2 WHERE id = $3`

	_, err := r.db.Exec(querry, newPassword, time.Now(), id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repository) GetCustomerByEmail(email string) (Customer, error) {
	querry := `SELECT * FROM customers WHERE email = $1`

	var customer Customer

	err := r.db.Get(&customer, querry, email)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (r *repository) DeleteCustomer(id uint) error {
	querry := `DELETE FROM customers WHERE id = $1`

	_, err := r.db.Exec(querry, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) IsEmailAvailable(email string) error {
	var id uint
	querry := `SELECT id FROM customers WHERE email = $1`
	err := r.db.QueryRowx(querry, email).Scan(&id)
	if err == sql.ErrNoRows && id == 0 {
		return nil
	}
	if err != nil {
		return err
	}

	return errors.New("email has been used by another")

}

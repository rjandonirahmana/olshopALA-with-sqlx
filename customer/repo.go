package customer

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

type Repository interface {
	RegisterUser(customer Customer) (Customer, error)
	UpdateCustomerPhone(id int, number string) error
	GetCustomerByID(id int) (Customer, error)
	ChangePassword(newPassword string, id int) error
	GetCustomerByEmail(email string) (Customer, error)
	ChangeAvatar(avatarFile string, id int) error
	DeleteCustomer(id int) error
}

func NewRepo(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) RegisterUser(customer Customer) (Customer, error) {

	querry := `INSERT INTO customers (name, phone, email, password, salt, avatar, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(querry, customer.Name, customer.Phone, customer.Email, customer.Password, customer.Salt, customer.Avatar, customer.CreatedAt, customer.UpdatedAt)

	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (r *repository) UpdateCustomerPhone(id int, number string) error {

	querry := `
	UPDATE 
		customers 
	SET 
		phone = ?, 
		updated_at = ?
	WHERE 
		id = ?
	`

	_, err := r.db.Exec(querry, number, time.Now(), id)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetCustomerByID(id int) (Customer, error) {
	querry := `SELECT * FROM customers WHERE id = ?`

	var customer Customer

	err := r.db.Get(&customer, querry, id)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (r *repository) ChangeAvatar(avatarFile string, id int) error {
	querry := `UPDATE customers SET avatar = ?, updated_at = ? WHERE id = ? `

	_, err := r.db.Exec(querry, avatarFile, time.Now(), id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repository) ChangePassword(newPassword string, id int) error {
	querry := `UPDATE customers SET password = ?, updated_at = ? WHERE id = ? `

	_, err := r.db.Exec(querry, newPassword, time.Now(), id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repository) GetCustomerByEmail(email string) (Customer, error) {
	querry := `SELECT * FROM customers WHERE email = ?`

	var customer Customer

	err := r.db.Get(&customer, querry, email)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (r *repository) DeleteCustomer(id int) error {
	querry := `DELETE FROM customers WHERE id = ?`

	_, err := r.db.Exec(querry, id)
	if err != nil {
		return err
	}

	return nil
}

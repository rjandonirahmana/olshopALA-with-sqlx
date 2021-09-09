package customer

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

type Repository interface {
	RegisterUser(customer Customer) (Customer, error)
	UpdateCustomerPhone(email string, number string) error
	GetCustomerByID(id int) (Customer, error)
	ChangePassword(newPassword string, id int) error
	GetCustomerByEmail(email string) (Customer, error)
	GetLastID() (int, error)
	ChangeAvatar(avatarFile string, id int) error
}

func NewRepo(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) RegisterUser(customer Customer) (Customer, error) {

	querry := `INSERT INTO customers (id, name, phone, email, password, salt, avatar, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(querry, customer.ID, customer.Name, customer.Phone, customer.Email, customer.Password, customer.Salt, customer.Avatar, customer.CreatedAt, customer.UpdatedAt)

	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (r *repository) UpdateCustomerPhone(email string, number string) error {

	querry := `
	UPDATE 
		customers 
	SET 
		phone = ? 
	WHERE 
		email = ?
	`

	_, err := r.db.Exec(querry, number, email)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetLastID() (int, error) {
	querry := `SELECT id FROM customers WHERE id = (SELECT MAX(id) FROM customers)`

	var value int
	err := r.db.Get(&value, querry)
	if err != nil {
		return 0, err
	}
	return value, nil

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
	querry := `UPDATE customers SET avatar = ? WHERE id = ? `

	_, err := r.db.Exec(querry, avatarFile, id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repository) ChangePassword(newPassword string, id int) error {
	querry := `UPDATE customers SET password = ? WHERE id = ? `

	_, err := r.db.Exec(querry, newPassword, id)

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

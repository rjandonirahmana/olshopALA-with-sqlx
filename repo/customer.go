package repo

import "github.com/jmoiron/sqlx"

type repository struct {
	db *sqlx.DB
}

type Repository interface {
	RegisterUser(customer Customer) error
	UpdateCustomerPhone(id int, customer Customer) error
	GetCustomerByID(id int) Customer
	ChangePassword(newPassword string, id int) error
}

func NewRepo(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) RegisterUser(customer Customer) error {

	querry := `INSERT INTO customers (id, name, phone, email, password, salt) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(querry, customer.ID, customer.Name, customer.Phone, customer.Email, customer.Password, "saltaja")

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateCustomerPhone(id int, customer Customer) error {

	querry := `
	UPDATE 
		customers 
	SET 
		phone = ? 
	WHERE 
		id = ?
	`

	_, err := r.db.Exec(querry, customer.Phone, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetCustomerByID(id int) Customer {
	querry := `SELECT * FROM customers WHERE id = ?`

	var customerdb CustomerDB

	err := r.db.Get(&customerdb, querry, id)
	if err != nil {
		return Customer{}
	}

	return Customer{
		Name:     customerdb.Name.String,
		ID:       int(customerdb.ID.Int32),
		Email:    customerdb.Email.String,
		Phone:    customerdb.Phone.Int64,
		Password: customerdb.Password.String,
		Salt:     customerdb.Salt.String,
	}
}

func (r *repository) ChangePassword(newPassword string, id int) error {
	querry := `UPDATE customers SET password = ? WHERE id = ? `

	_, err := r.db.Exec(querry, newPassword, id)

	if err != nil {
		return err
	}

	return nil
}

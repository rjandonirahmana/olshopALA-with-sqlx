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
	UpdateCustomerPhone(email string, number string) error
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
	rows, _ := r.db.Query("show columns from customers")
	cols := []string{}
	for rows.Next() {
		var (
			Field   string
			Type    string
			Null    string
			Key     string
			Default string
			Extra   string
		)
		rows.Scan(&Field, &Type, &Null, &Key, &Default, &Extra)
		cols = append(cols, Field)
	}
	colString := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s", cols[0], cols[1], cols[2], cols[3], cols[4], cols[5], cols[6], cols[7], cols[8])
	fmt.Println(colString)
	querry := `INSERT INTO customers (` + colString + `) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(querry, customer.ID, customer.Name, customer.Phone, customer.Email, customer.Password, customer.Salt, customer.Avatar, customer.CreatedAt, customer.UpdatedAt)

	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (r *repository) UpdateCustomerPhone(email string, number string) error {
	rows, _ := r.db.Query("show columns from customers")
	cols := []string{}
	for rows.Next() {
		var (
			Field   string
			Type    string
			Null    string
			Key     string
			Default string
			Extra   string
		)
		rows.Scan(&Field, &Type, &Null, &Key, &Default, &Extra)
		cols = append(cols, Field)
	}
	// colString := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,", cols[0], cols[1], cols[2],cols[3], cols[4], cols[5],cols[6], cols[7], cols[8],)
	// querry := `
	// UPDATE
	// 	customers
	// SET
	// 	phone = ?,
	// 	updated_at = ?
	// WHERE
	// 	email = ?
	// `
	querry := fmt.Sprintf("UPDATE customers SET %s = ?, %s = ? WHERE %s = ?", cols[2], cols[8], cols[3])
	_, err := r.db.Exec(querry, number, time.Now(), email)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetCustomerByID(id int) (Customer, error) {
	rows, _ := r.db.Query("show columns from customers")
	cols := []string{}
	for rows.Next() {
		var (
			Field   string
			Type    string
			Null    string
			Key     string
			Default string
			Extra   string
		)
		rows.Scan(&Field, &Type, &Null, &Key, &Default, &Extra)
		cols = append(cols, Field)
	}
	// colString := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,", cols[0], cols[1], cols[2],cols[3], cols[4], cols[5],cols[6], cols[7], cols[8],)
	querry := `SELECT * FROM customers WHERE ` + cols[0] + ` = ?`

	var customer Customer

	err := r.db.Get(&customer, querry, id)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (r *repository) ChangeAvatar(avatarFile string, id int) error {
	rows, _ := r.db.Query("show columns from customers")
	cols := []string{}
	for rows.Next() {
		var (
			Field   string
			Type    string
			Null    string
			Key     string
			Default string
			Extra   string
		)
		rows.Scan(&Field, &Type, &Null, &Key, &Default, &Extra)
		cols = append(cols, Field)
	}
	// colString := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,", cols[0], cols[1], cols[2],cols[3], cols[4], cols[5],cols[6], cols[7], cols[8],)
	// querry := `
	// UPDATE
	// 	customers
	// SET
	// 	phone = ?,
	// 	updated_at = ?
	// WHERE
	// 	email = ?
	// `
	querry := fmt.Sprintf("UPDATE customers SET %s = ?, %s = ? WHERE %s = ?", cols[6], cols[8], cols[0])
	// querry := `UPDATE customers SET avatar = ?, updated_at = ? WHERE id = ? `

	_, err := r.db.Exec(querry, avatarFile, time.Now(), id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repository) ChangePassword(newPassword string, id int) error {
	rows, _ := r.db.Query("show columns from customers")
	cols := []string{}
	for rows.Next() {
		var (
			Field   string
			Type    string
			Null    string
			Key     string
			Default string
			Extra   string
		)
		rows.Scan(&Field, &Type, &Null, &Key, &Default, &Extra)
		cols = append(cols, Field)
	}
	// colString := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,", cols[0], cols[1], cols[2],cols[3], cols[4], cols[5],cols[6], cols[7], cols[8],)
	// querry := `
	// UPDATE
	// 	customers
	// SET
	// 	phone = ?,
	// 	updated_at = ?
	// WHERE
	// 	email = ?
	// `
	querry := fmt.Sprintf("UPDATE customers SET %s = ?, %s = ? WHERE %s = ?", cols[4], cols[8], cols[0])
	// querry := `UPDATE customers SET password = ?, updated_at = ? WHERE id = ? `

	_, err := r.db.Exec(querry, newPassword, time.Now(), id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repository) GetCustomerByEmail(email string) (Customer, error) {
	rows, _ := r.db.Query("show columns from customers")
	cols := []string{}
	for rows.Next() {
		var (
			Field   string
			Type    string
			Null    string
			Key     string
			Default string
			Extra   string
		)
		rows.Scan(&Field, &Type, &Null, &Key, &Default, &Extra)
		cols = append(cols, Field)
	}
	// colString := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,", cols[0], cols[1], cols[2],cols[3], cols[4], cols[5],cols[6], cols[7], cols[8],)
	querry := `SELECT * FROM customers WHERE ` + cols[3] + ` = ?`

	var customer Customer
	row, _ := r.db.Query(querry, email)
	for row.Next() {
		var (
			ID        int
			Name      string
			Email     string
			Phone     string
			Password  string
			Salt      string
			Avatar    string
			CreatedAt string
			UpdatedAt string
		)
		row.Scan(&ID, &Name, &Email, &Phone, &Password, &Salt, &Avatar, &CreatedAt, &UpdatedAt)
		fmt.Println(ID, Name, Email, Phone)
		break
	}
	err := r.db.Get(&customer, querry, email)
	fmt.Println(err)
	if err != nil {
		fmt.Println("awww")
		return Customer{}, err
	}

	return customer, nil
}

func (r *repository) DeleteCustomer(id int) error {
	rows, _ := r.db.Query("show columns from customers")
	cols := []string{}
	for rows.Next() {
		var (
			Field   string
			Type    string
			Null    string
			Key     string
			Default string
			Extra   string
		)
		rows.Scan(&Field, &Type, &Null, &Key, &Default, &Extra)
		cols = append(cols, Field)
	}
	// colString := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,", cols[0], cols[1], cols[2],cols[3], cols[4], cols[5],cols[6], cols[7], cols[8],)
	querry := `DELETE FROM customers WHERE ` + cols[0] + ` = ?`

	_, err := r.db.Exec(querry, id)
	if err != nil {
		return err
	}

	return nil
}

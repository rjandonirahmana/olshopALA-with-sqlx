package payment

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func (r *Repo) AddPaymenthMethod(payment Payment) error {
	querry := `INSERT INTO payment (id, methode) VALUES (?, ?)`

	_, err := r.db.Exec(querry, payment.ID, payment.Methode)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

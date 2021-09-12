package transaction

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type repoTransaction struct {
	db *sqlx.DB
}

type RepoTransaction interface {
	SumPriceBoughtById(id int) int
	GetDetailTransaction(id int) (Transactions, error)
	InserTransaction(t Transactions) error
	GetLastTransaction() (int, error)
}

func NewTransactionRepo(db *sqlx.DB) *repoTransaction {
	return &repoTransaction{db: db}
}

func (r *repoTransaction) SumPriceBoughtById(id int) int {
	querry := `SELECT SUM(price) FROM transactions WHERE customer_id = ? `

	var total int

	err := r.db.Get(&total, querry, id)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	return total

}

func (r *repoTransaction) GetDetailTransaction(id int) (Transactions, error) {
	querry := `SELECT * FROM transactions WHERE id = ?`
	var transaction Transactions
	err := r.db.Get(&transaction, querry, id)

	if err != nil {
		fmt.Println(err)
		return Transactions{}, err
	}

	return transaction, nil
}

func (r *repoTransaction) InserTransaction(t Transactions) error {

	querry := `INSERT INTO transactions
	(id, customer_id, price, created_at, max_time, shopcart_id, payment_id)
	VALUES(?, ?, ?, ?, ?, ?, ?);`

	_, err := r.db.Exec(querry, t.ID, t.CustomerID, t.Price, t.CreatedAt, t.MaxTime, t.ShopCartID, t.PaymentID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repoTransaction) GetLastTransaction() (int, error) {
	querry := `SELECT id FROM transactions WHERE id = (SELECT MAX(id) FROM transactions`

	var value int
	err := r.db.Get(&value, querry)
	if value == 0 {
		return value + 1, nil
	}
	if err != nil {
		return 0, err
	}
	return value + 1, nil

}

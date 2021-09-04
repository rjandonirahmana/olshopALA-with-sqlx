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
	GetDetailProductByID(id int) Transactions
}

func NewTransaction(db *sqlx.DB) *repoTransaction {
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

func (r *repoTransaction) GetDetailProductByID(id int) Transactions {
	querry := `SELECT * FROM transactions WHERE id = ?`
	var transaction Transactions
	err := r.db.Get(&transaction, querry, id)

	if err != nil {
		fmt.Println(err)
		return Transactions{}
	}

	return transaction
}

func (r *repoTransaction) InserTransaction(t Transactions) error {

	querry := `INSERT INTO transactions
	(id, product_id, customer_id, quantity, price, created_at)
	VALUES(?,?,?,?,?,?)`

	_, err := r.db.Exec(querry, t.ID, t.ID_product, t.CustomerID, t.Quantity, t.Price, t.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

package transaction

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type repoTransaction struct {
	db *sqlx.DB
}

type RepoTransaction interface {
	GetDetailTransaction(id int) (Transactions, error)
	InserTransaction(t Transactions) error
	CheckTransaction(cartid int) (int, error)
}

func NewTransactionRepo(db *sqlx.DB) *repoTransaction {
	return &repoTransaction{db: db}
}

func (r *repoTransaction) GetDetailTransaction(id int) (Transactions, error) {
	querry := `SELECT * FROM transactions WHERE id = ?`
	var trans Transactions
	err := r.db.Get(&trans, querry, id)

	if err != nil {
		return Transactions{}, err
	}

	return trans, nil
}

func (r *repoTransaction) InserTransaction(t Transactions) error {

	querry := `INSERT INTO transactions
	(id, customer_id, price, created_at, max_time, shopcart_id, payment_id)
	VALUES(?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(querry, t.ID, t.CustomerID, t.Price, t.CreatedAt, t.MaxTime, t.ShopCartID, t.PaymentID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repoTransaction) CheckTransaction(cartid int) (int, error) {
	querry := `SELECT id FROM transactions WHERE shopcart_id = ? AND max_time > ?`

	//select id transacttion biar di passs ke detail transaction kalo ada

	row, err := r.db.Queryx(querry, cartid, time.Now())
	var value int

	for row.Next() {
		err = row.Scan(&value)
		if err != nil {
			return 0, err
		}
	}
	if value == 0 {
		return 0, nil
	}
	defer row.Close()

	if row.Err() != nil {
		return 0, err
	}

	return value, nil
}

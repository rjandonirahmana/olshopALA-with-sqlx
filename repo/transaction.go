package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type repoTransaction struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) *repoTransaction {
	return &repoTransaction{db: db}
}

func (r *repoTransaction) SumPriceBoughtById(id int) int {
	querry := `SELECT SUM(price) FROM transactions WHERE costumer_id = ? `

	var total []int

	err := r.db.Select(&total, querry, id)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	sum := 0
	for _, v := range total {
		sum += v
	}

	return sum
}

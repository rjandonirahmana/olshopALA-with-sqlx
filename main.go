package main

import (
	"fmt"
	"graphql/transaction"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("", ":@(localhost:3306)/")
	if err != nil {
		log.Fatalln(err)
	}

	product := transaction.NewTransaction(db)

	coba := product.SumPriceBoughtById(1)
	fmt.Println(coba)
}

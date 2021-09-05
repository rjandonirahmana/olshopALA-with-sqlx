package main

import (
	"graphql/customer"
	"graphql/handler"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("mysql", ":@(localhost:3306)/?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	customerdb := customer.NewRepo(db)
	customerserv := customer.NewCustomerService(customerdb)
	customerHandler := handler.NewHandlerCustomer(customerserv)

	c := gin.New()
	c.POST("/register", customerHandler.CreateCustomer)
	c.POST("/login", customerHandler.Login)

	c.Run(":8080")
}

package main

import (
	"graphql/handler"
	"graphql/repo"
	"graphql/usecase"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("mysql", ":@(localhost:3306)/")
	if err != nil {
		log.Fatalln(err)
	}

	repoUser := repo.NewRepo(db)
	usecaseUser := usecase.NewCustomerService(repoUser)
	hanlderUser := handler.NewHandlerCustomer(usecaseUser)

	c := gin.Default()

	c.POST("/customer", hanlderUser.CreateCustomer)

}

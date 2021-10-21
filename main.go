package main

import (
	"fmt"
	"log"
	"olshop/auth"
	"olshop/customer"
	"olshop/handler"
	c "olshop/handler/customer"
	p "olshop/handler/product"
	s "olshop/handler/seller"
	"olshop/product"
	"olshop/seller"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("disni??")
		log.Fatal("Error loading .env file")
	}
	dbUserName := os.Getenv("DB_username")
	dbName := os.Getenv("DB_name")
	dbPass := os.Getenv("DB_password")
	secretkey1 := os.Getenv("SECRET_KEY")
	secretkey2 := os.Getenv("SECRET_KEY2")
	secetpasscustomer := os.Getenv("SECRETPASS")
	secretpassseller := os.Getenv("SECRETPASS2")

	dbString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", dbUserName, dbPass, dbName)
	db, err := sqlx.Connect("pgx", dbString)

	if err != nil {
		panic(err)
	}

	auth := auth.NewService(secretkey1, secretkey2)
	customerdb := customer.NewRepo(db)
	productdb := product.NewRepoProduct(db)
	sellerdb := seller.NewRepository(db)
	// product, err := productdb.SearchAndByorder("iphone", 0, 0)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(product)

	// transactiondb := transaction.NewTransactionRepo(db)

	authMiddleWare := handler.NewMiddleWare()
	customerserv := customer.NewCustomerService(customerdb, secetpasscustomer)
	productServ := product.NewService(productdb)
	sellerserv := seller.NewService(sellerdb, secretpassseller)
	// // transactionServ := transaction.NewTransactionService(transactiondb, productdb)

	productHanlder := p.NewProductHandler(productServ)
	customerHandler := c.NewHandlerCustomer(customerserv, auth)
	sellerHandler := s.NewHandlerSeller(sellerserv, auth)
	// // transactionHandler := handler.NewTransactionHandler(transactionServ)

	c := gin.Default()
	api := c.Group("/api/v1")

	//producr
	api.GET("/category", productHanlder.GetProductByCategory)
	api.GET("/product/:id", productHanlder.GetProductByID)

	//customer
	api.POST("/register", customerHandler.CreateCustomer)
	api.POST("/login", customerHandler.Login)
	api.PUT("/phone", authMiddleWare.AuthMiddleWareCustomer(auth, customerserv), customerHandler.UpdatePhoneCustomer)
	// // api.PUT("/avatar", authMiddleWare(auth, customerserv), customerHandler.UpdateAvatar)
	// // api.PUT("password", authMiddleWare(auth, customerserv), customerHandler.UpdatePassword)
	// // api.DELETE("/account", authMiddleWare(auth, customerserv), customerHandler.DeleteAccount)
	// // api.POST("/addcart", authMiddleWare(auth, customerserv), productHanlder.CreateShopCart)

	//seller
	api.POST("/registerseller", sellerHandler.Register)
	api.POST("/loginseller", sellerHandler.Login)
	// // api.GET("/listshopcart", authMiddleWare(auth, customerserv), productHanlder.GetListProductShopCart)
	// // api.GET("/shopcartcustomer", authMiddleWare(auth, customerserv), productHanlder.GetAllCartCustomer)
	// // api.PUT("/decreaseproduct", authMiddleWare(auth, customerserv), productHanlder.DecreaseQuantity)
	// // api.DELETE("/productshopcart", authMiddleWare(auth, customerserv), productHanlder.DeleteProductShopcart)
	// // api.POST("/transaction", authMiddleWare(auth, customerserv), transactionHandler.CreateTransaction)

	c.Run(":8080")
}

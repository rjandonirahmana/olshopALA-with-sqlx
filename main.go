package main

import (
	"fmt"
	"log"
	"net/http"
	"olshop/auth"
	"olshop/customer"
	"olshop/handler"
	"olshop/product"
	"olshop/seller"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbname := os.Getenv("DB_name")
	dbpassword := os.Getenv("DB_password")
	dbuser := os.Getenv("DB_username")
	// secretkey1 := os.Getenv("SECRET_KEY")
	// secretkey2 := os.Getenv("SECRET_KEY2")

	conection_db := fmt.Sprintf("%s:%s@(localhost:3306)/%s?parseTime=true", dbuser, dbpassword, dbname)

	db, err := sqlx.Connect("mysql", conection_db)
	if err != nil {
		log.Fatalln(err)
	}

	// auth := auth.NewService(secretkey1, secretkey2)
	// customerdb := customer.NewRepo(db)
	productdb := product.NewRepoProduct(db)
	product, err := productdb.SearchAndByorder("iphone", 0, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(product)

	// transactiondb := transaction.NewTransactionRepo(db)

	// customerserv := customer.NewCustomerService(customerdb)
	// productServ := product.NewService(productdb)
	// // transactionServ := transaction.NewTransactionService(transactiondb, productdb)

	// productHanlder := handler.NewProductHandler(productServ)
	// // customerHandler := handler.NewHandlerCustomer(customerserv, auth)
	// // transactionHandler := handler.NewTransactionHandler(transactionServ)

	// c := gin.Default()
	// api := c.Group("/api/v1")

	// api.GET("/category/:id", productHanlder.GetProductByCategory)
	// api.GET("/product/:id", productHanlder.GetProductByID)
	// // api.POST("/register", customerHandler.CreateCustomer)
	// // api.POST("/login", customerHandler.Login)
	// // api.PUT("/phone", authMiddleWare(auth, customerserv), customerHandler.UpdatePhoneCustomer)
	// // api.PUT("/avatar", authMiddleWare(auth, customerserv), customerHandler.UpdateAvatar)
	// // api.PUT("password", authMiddleWare(auth, customerserv), customerHandler.UpdatePassword)
	// // api.DELETE("/account", authMiddleWare(auth, customerserv), customerHandler.DeleteAccount)
	// // api.POST("/addcart", authMiddleWare(auth, customerserv), productHanlder.CreateShopCart)

	// // api.GET("/listshopcart", authMiddleWare(auth, customerserv), productHanlder.GetListProductShopCart)
	// // api.GET("/shopcartcustomer", authMiddleWare(auth, customerserv), productHanlder.GetAllCartCustomer)
	// // api.PUT("/decreaseproduct", authMiddleWare(auth, customerserv), productHanlder.DecreaseQuantity)
	// // api.DELETE("/productshopcart", authMiddleWare(auth, customerserv), productHanlder.DeleteProductShopcart)
	// // api.POST("/transaction", authMiddleWare(auth, customerserv), transactionHandler.CreateTransaction)

	// c.Run(":8080")
}

func authMiddleWareCustomer(auth auth.Service, service customer.CustomerInt) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := handler.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		var tokenString string
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := auth.ValidateToken(tokenString)
		if err != nil {
			response := handler.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := handler.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		customerID := int(claim["customer_id"].(float64))

		customer, err := service.GetCustomerByID(customerID)
		if err != nil {
			response := handler.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentCustomer", customer)
	}
}

func authMiddleWareSeller(auth auth.Service, service seller.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := handler.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		var tokenString string
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := auth.ValidateTokenSeller(tokenString)
		if err != nil {
			response := handler.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := handler.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		sellerID := int(claim["seller_id"].(float64))

		seller, err := service.GetSellerByID(sellerID)
		if err != nil {
			response := handler.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("seller", seller)
	}
}

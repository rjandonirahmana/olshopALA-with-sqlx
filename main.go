package main

import (
	"log"
	"net/http"
	"olshop/auth"
	"olshop/customer"
	"olshop/handler"
	"olshop/product"
	"olshop/transaction"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("mysql", "root:12345@(localhost:3306)/olshopALA?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	transactiondb := transaction.NewTransactionRepo(db)

	customerserv := customer.NewCustomerService(customerdb)
	productServ := product.NewService(productdb)
	transactionServ := transaction.NewTransactionService(transactiondb, productdb)

	productHanlder := handler.NewProductHandler(productServ)
	customerHandler := handler.NewHandlerCustomer(customerserv, auth)
	transactionHandler := handler.NewTransactionHandler(transactionServ)

	c := gin.Default()
	api := c.Group("/api/v1")

	api.GET("/productcategory", productHanlder.GetProductByCategory)
	api.POST("/register", customerHandler.CreateCustomer)
	api.POST("/login", customerHandler.Login)
	api.PUT("/phone", authMiddleWare(auth, customerserv), customerHandler.UpdatePhoneCustomer)
	api.PUT("/avatar", authMiddleWare(auth, customerserv), customerHandler.UpdateAvatar)
	api.PUT("password", authMiddleWare(auth, customerserv), customerHandler.UpdatePassword)
	api.DELETE("/account", authMiddleWare(auth, customerserv), customerHandler.DeleteAccount)
	api.POST("/addcart", authMiddleWare(auth, customerserv), productHanlder.CreateShopCart)

	api.POST("/insertshopcart", authMiddleWare(auth, customerserv), productHanlder.InsertToShopCart)
	api.GET("/listshopcart", authMiddleWare(auth, customerserv), productHanlder.GetListProductShopCart)
	api.GET("/shopcartcustomer", authMiddleWare(auth, customerserv), productHanlder.GetAllCartCustomer)
	api.PUT("/decreaseproduct", authMiddleWare(auth, customerserv), productHanlder.DecreaseQuantity)
	api.DELETE("/productshopcart", authMiddleWare(auth, customerserv), productHanlder.DeleteProductShopcart)
	api.POST("/transaction", authMiddleWare(auth, customerserv), transactionHandler.CreateTransaction)

	c.Run(":8080")
}

func authMiddleWare(auth auth.Service, service customer.CustomerInt) gin.HandlerFunc {
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

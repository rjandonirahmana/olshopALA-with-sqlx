package main

import (
	"graphql/auth"
	"graphql/customer"
	"graphql/handler"
	"graphql/product"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("mysql", ":@(localhost:3306)/?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	auth := auth.NewService()
	customerdb := customer.NewRepo(db)
	productdb := product.NewRepoProduct(db)
	customerserv := customer.NewCustomerService(customerdb)
	productServ := product.NewService(productdb)

	productHanlder := handler.NewProductHandler(productServ)
	customerHandler := handler.NewHandlerCustomer(customerserv, auth)

	c := gin.New()

	c.GET("productCategory", productHanlder.GetProductByCategory)
	c.POST("/register", customerHandler.CreateCustomer)
	c.POST("/login", customerHandler.Login)
	c.PUT("/phone", authMiddleWare(auth, customerserv), customerHandler.UpdatePhoneCustomer)
	c.PUT("/avatar", authMiddleWare(auth, customerserv), customerHandler.UpdateAvatar)
	c.POST("/addcart", authMiddleWare(auth, customerserv), productHanlder.CreateShopCart)
	c.POST("/addshopcart", authMiddleWare(auth, customerserv), productHanlder.InsertToShopCart)
	c.GET("/listshopcart", authMiddleWare(auth, customerserv), productHanlder.GetListProductShopCart)
	c.DELETE("/listshopcart", authMiddleWare(auth, customerserv), productHanlder.DeleteProductShopcart)

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

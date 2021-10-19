package handler

import (
	"net/http"
	"olshop/auth"
	"olshop/customer"
	"olshop/seller"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type middleWare struct{}

func NewMiddleWare() *middleWare {
	return &middleWare{}
}

func (m *middleWare) AuthMiddleWareCustomer(auth auth.Service, service customer.CustomerInt) gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, err := c.Request.Cookie("tokencustomer")
		if err != nil {
			response := APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		token, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			response := APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		customerID := int(claim["customer_id"].(float64))

		customer, err := service.GetCustomerByID(customerID)
		if err != nil {
			response := APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentCustomer", customer)
	}
}

func (m *middleWare) AuthMiddleWareSeller(auth auth.Service, service seller.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, err := c.Request.Cookie("tokenseller")
		if err != nil {
			response := APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		token, err := auth.ValidateTokenSeller(cookie.Value)
		if err != nil {
			response := APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		sellerID := int(claim["seller_id"].(float64))

		seller, err := service.GetSellerByID(sellerID)
		if err != nil {
			response := APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("CurrentSeller", seller)
	}
}

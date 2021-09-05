package handler

import (
	"errors"
	"graphql/customer"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handlerCustomer struct {
	usecase customer.CustomerInt
}

func NewHandlerCustomer(use customer.CustomerInt) *handlerCustomer {
	return &handlerCustomer{usecase: use}
}

func (h *handlerCustomer) CreateCustomer(c *gin.Context) {
	var customers customer.InputCustomer

	c.ShouldBindJSON(&customers)
	if customers.Password != customers.ConfirmPassword {
		response := APIResponse("password and confirm password is different", http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	customerSave := customer.Customer{}
	customerSave.Name = customers.Name
	customerSave.Email = customers.Email
	customerSave.Password = customers.Password

	customer, err := h.usecase.Register(customerSave)

	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("failed to create customer"))
		return
	}

	response := APIResponse("account successfully created", http.StatusOK, "success", customer)
	c.JSON(http.StatusUnprocessableEntity, response)

}

func (h *handlerCustomer) Login(c *gin.Context) {
	var input customer.InputLogin

	c.ShouldBindJSON(&input)

	customer, err := h.usecase.LoginCustomer(input)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusForbidden, "failed", err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("success login", 200, "success", customer)

	c.JSON(http.StatusOK, response)

}

// func (h *handlerCustomer) UpdatePhone(c *gin.Context) {
// 	phone, err := strconv.ParseInt(c.Request.FormValue("phone"), 10, 10)

// }

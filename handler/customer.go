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
		c.JSON(http.StatusBadRequest, errors.New("password and confirmed password is different"))
	}

	customerSave := customer.Customer{}
	customerSave.Name = customers.Name
	customerSave.Email = customers.Email
	customerSave.Password = customers.Password

	err := h.usecase.Register(customerSave)

	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("failed to create customer"))
		return
	}

	c.JSON(http.StatusOK, nil)

}

func (h *handlerCustomer) Login(c *gin.Context) {
	var input customer.InputLogin

	c.ShouldBindJSON(&input)

	customer, err := h.usecase.LoginCustomer(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("failed to login customer"))
		return
	}

	response := APIResponse("success login", 200, "success", customer)

	c.JSON(http.StatusOK, response)

}

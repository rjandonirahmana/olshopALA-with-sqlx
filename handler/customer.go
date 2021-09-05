package handler

import (
	"errors"
	"fmt"
	"graphql/auth"
	"graphql/customer"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type handlerCustomer struct {
	usecase customer.CustomerInt
	auth    auth.Service
}

func NewHandlerCustomer(use customer.CustomerInt, auth auth.Service) *handlerCustomer {
	return &handlerCustomer{usecase: use, auth: auth}
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

	token, err := h.auth.GenerateToken(customer.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("failed to generate token"))
		return

	}

	response := APIResponse("account successfully created", http.StatusOK, fmt.Sprintf("success created token %s", token), customer)
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

	token, err := h.auth.GenerateToken(customer.ID)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusForbidden, "failed", err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("success login", 200, fmt.Sprintf("success create token %s", token), customer)

	c.JSON(http.StatusOK, response)

}

func (h *handlerCustomer) UpdatePhoneCustomer(c *gin.Context) {
	phone, err := strconv.ParseInt(c.Request.FormValue("phone"), 10, 64)

	if err != nil {
		response := APIResponse(err.Error(), http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentCustomer := c.MustGet("currentCustomer").(customer.Customer)
	customerEmail := currentCustomer.Email

	err = h.usecase.UpdateCustomerPhone(phone, customerEmail)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updatedCustomer, err := h.usecase.GetCustomerByID(currentCustomer.ID)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("success login", 200, fmt.Sprintf("successfully update customer number %d", phone), updatedCustomer)

	c.JSON(http.StatusOK, response)

}

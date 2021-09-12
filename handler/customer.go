package handler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"olshop/auth"
	"olshop/customer"

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
		respones := APIResponse("failed to create account", http.StatusOK, fmt.Sprintf("%v", err.Error()), nil)
		c.JSON(http.StatusBadRequest, respones)
		return
	}

	token, err := h.auth.GenerateToken(customer.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("failed to generate token"))
		return

	}

	response := APIResponse("account successfully created", http.StatusOK, fmt.Sprintf("new customer suceessfully created with token %s", token), customer)
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

	response := APIResponse("success login", 200, fmt.Sprintf("%s's new temporary token %s", customer.Email, token), customer)

	c.JSON(http.StatusOK, response)

}

func (h *handlerCustomer) UpdatePhoneCustomer(c *gin.Context) {
	phone := c.Request.FormValue("phone")

	currentCustomer := c.MustGet("currentCustomer").(customer.Customer)
	customerEmail := currentCustomer.Email

	err := h.usecase.UpdateCustomerPhone(phone, customerEmail)
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

	response := APIResponse("success", 200, fmt.Sprintf("%s's number has been updated successfully ", currentCustomer.Email), updatedCustomer)

	c.JSON(http.StatusOK, response)

}

func (h *handlerCustomer) UpdateAvatar(c *gin.Context) {
	avatar, foto, err := c.Request.FormFile("avatar")
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentCustomer := c.MustGet("currentCustomer").(customer.Customer)

	file, err := ioutil.ReadAll(avatar)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = h.usecase.ChangeProfile(file, currentCustomer.Email, currentCustomer.ID)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	path := fmt.Sprintf("images/%d-%s", currentCustomer.ID, currentCustomer.Email)

	err = c.SaveUploadedFile(foto, path)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("avatar", 200, fmt.Sprintf("%s's avatar has successfuly been updated", currentCustomer.Email), currentCustomer)

	c.JSON(http.StatusOK, response)
}

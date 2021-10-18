package handler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"olshop/auth"
	"olshop/customer"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	err := c.ShouldBindJSON(&customers)

	if err != nil {
		response := APIResponse("failed", http.StatusUnprocessableEntity, "make sure to input whole field correctly", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	validate := validator.New()
	err = validate.Struct(&customers)
	if err != nil {
		response := APIResponse("failed", http.StatusForbidden, "make sure to input whole field correctly", nil)
		c.JSON(http.StatusForbidden, response)
		return
	}

	customerSave := customer.Customer{}
	customerSave.Name = customers.Name
	customerSave.Email = customers.Email
	customerSave.Password = customers.Password

	customer, err := h.usecase.Register(customerSave)

	if err != nil {
		respones := APIResponse("failed", http.StatusForbidden, fmt.Sprintf("%v", err.Error()), nil)
		c.JSON(http.StatusForbidden, respones)
		return
	}

	token, err := h.auth.GenerateToken(customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.New("failed to generate token"))
		return

	}

	response := ResponseAPIToken("success", http.StatusOK, fmt.Sprintf("new customer successfully created with id %d", customer.ID), customer, token)
	c.JSON(http.StatusOK, response)

}

func (h *handlerCustomer) Login(c *gin.Context) {
	var input customer.InputLogin

	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := APIResponse("failed", http.StatusForbidden, err.Error(), err)
		c.JSON(http.StatusForbidden, response)
		return
	}

	customer, err := h.usecase.LoginCustomer(input)
	if err != nil {
		response := APIResponse("failed", http.StatusForbidden, err.Error(), err)
		c.JSON(http.StatusForbidden, response)
		return
	}

	token, err := h.auth.GenerateToken(customer.ID)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := ResponseAPIToken("success", http.StatusOK, fmt.Sprintf("%s's account login successfully ", customer.Email), customer, token)

	c.JSON(http.StatusOK, response)

}

func (h *handlerCustomer) UpdatePhoneCustomer(c *gin.Context) {
	phone := c.Request.FormValue("phone")

	currentCustomer := c.MustGet("currentCustomer").(customer.Customer)

	err := h.usecase.UpdateCustomerPhone(phone, currentCustomer.ID)
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

	currentCustomer, err = h.usecase.ChangeProfile(file, currentCustomer.Email, currentCustomer.ID)
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

func (h *handlerCustomer) UpdatePassword(c *gin.Context) {
	password := c.Request.FormValue("password")
	newPassword := c.Request.FormValue("newpassword")

	customer := c.MustGet("currentCustomer").(customer.Customer)

	customer, err := h.usecase.ChangePassword(password, newPassword, customer.ID)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	response := APIResponse("success", 200, fmt.Sprintf("%s's password has successfuly been updated", customer.Email), customer)

	c.JSON(http.StatusOK, response)

}

func (h *handlerCustomer) DeleteAccount(c *gin.Context) {
	password := c.Request.FormValue("password")

	customer := c.MustGet("currentCustomer").(customer.Customer)

	err := h.usecase.DeleteCustomer(customer.ID, password)

	if err != nil {
		response := APIResponse(err.Error(), http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("success", 200, fmt.Sprintf("%s's account has been deleted", customer.Email), nil)

	c.JSON(http.StatusOK, response)

}

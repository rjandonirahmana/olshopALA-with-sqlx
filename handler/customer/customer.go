package customer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"olshop/auth"
	"olshop/customer"
	"olshop/handler"
	"time"

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
		response := handler.APIResponse("make sure to input whole field correctly", http.StatusUnprocessableEntity, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	validate := validator.New()
	err = validate.Struct(&customers)
	if err != nil {
		response := handler.APIResponse("your password need to be 8 length and make sure your confirm password match your password", http.StatusForbidden, "failed", err.Error())
		c.JSON(http.StatusForbidden, response)
		return
	}

	customerSave := customer.Customer{}
	customerSave.Name = customers.Name
	customerSave.Email = customers.Email
	customerSave.Password = customers.Password

	customer, err := h.usecase.Register(customerSave)

	if err != nil {
		respones := handler.APIResponse(fmt.Sprintf("%v", err.Error()), http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusForbidden, respones)
		return
	}

	token, err := h.auth.GenerateToken(customer.ID)
	if err != nil {
		respones := handler.APIResponse("failed to generate token", http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusForbidden, respones)
		return

	}

	cookie := http.Cookie{
		Name:    "tokencustomer",
		Value:   token,
		Expires: time.Now().Add(12 * time.Hour),
	}

	http.SetCookie(c.Writer, &cookie)

	response := handler.ResponseAPIToken("new customer successfully created with id", http.StatusOK, "suceess", customer, token)
	c.JSON(http.StatusOK, response)

}

func (h *handlerCustomer) Login(c *gin.Context) {
	var input customer.InputLogin

	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := handler.APIResponse("make sure to input field correctly", http.StatusForbidden, "failed", err.Error())
		c.JSON(http.StatusForbidden, response)
		return
	}

	customer, err := h.usecase.LoginCustomer(input)
	if err != nil {
		response := handler.APIResponse("failed to login", http.StatusForbidden, "failed", err.Error())
		c.JSON(http.StatusForbidden, response)
		return
	}

	token, err := h.auth.GenerateToken(customer.ID)
	if err != nil {
		response := handler.APIResponse("failed to generate token", http.StatusInternalServerError, "failed", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := handler.ResponseAPIToken(fmt.Sprintf("%s's account login successfully ", customer.Email), http.StatusOK, "success", customer, token)

	cookie := http.Cookie{
		Name:    "tokencustomer",
		Value:   token,
		Expires: time.Now().Add(12 * time.Hour),
	}

	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusOK, response)

}

func (h *handlerCustomer) UpdatePhoneCustomer(c *gin.Context) {
	phone := c.Request.FormValue("phone")

	currentCustomer := c.MustGet("currentCustomer").(customer.Customer)

	err := h.usecase.UpdateCustomerPhone(phone, currentCustomer.ID)
	if err != nil {
		response := handler.APIResponse("failed to update phone number", http.StatusForbidden, "failed", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updatedCustomer, err := h.usecase.GetCustomerByID(currentCustomer.ID)
	if err != nil {
		response := handler.APIResponse("failed to get customer", http.StatusForbidden, "failed", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := handler.APIResponse(fmt.Sprintf("%s's number has been updated successfully ", currentCustomer.Email), 200, "success", updatedCustomer)

	c.JSON(http.StatusOK, response)

}

func (h *handlerCustomer) UpdateAvatar(c *gin.Context) {
	avatar, foto, err := c.Request.FormFile("avatar")
	if err != nil {
		response := handler.APIResponse("failed to get file", http.StatusUnprocessableEntity, "failed", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentCustomer := c.MustGet("currentCustomer").(customer.Customer)

	file, err := ioutil.ReadAll(avatar)
	if err != nil {
		response := handler.APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentCustomer, err = h.usecase.ChangeProfile(file, currentCustomer.Email, currentCustomer.ID)
	if err != nil {
		response := handler.APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	path := fmt.Sprintf("images/%d-%s", currentCustomer.ID, currentCustomer.Email)

	err = c.SaveUploadedFile(foto, path)
	if err != nil {
		response := handler.APIResponse(err.Error(), http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := handler.APIResponse("avatar", 200, fmt.Sprintf("%s's avatar has successfuly been updated", currentCustomer.Email), currentCustomer)

	c.JSON(http.StatusOK, response)
}

func (h *handlerCustomer) UpdatePassword(c *gin.Context) {
	password := c.Request.FormValue("password")
	newPassword := c.Request.FormValue("newpassword")

	customer := c.MustGet("currentCustomer").(customer.Customer)

	customer, err := h.usecase.ChangePassword(password, newPassword, customer.ID)
	if err != nil {
		response := handler.APIResponse(err.Error(), http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	response := handler.APIResponse("success", 200, fmt.Sprintf("%s's password has successfuly been updated", customer.Email), customer)

	c.JSON(http.StatusOK, response)

}

func (h *handlerCustomer) DeleteAccount(c *gin.Context) {
	password := c.Request.FormValue("password")

	customer := c.MustGet("currentCustomer").(customer.Customer)

	err := h.usecase.DeleteCustomer(customer.ID, password)

	if err != nil {
		response := handler.APIResponse(err.Error(), http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := handler.APIResponse("success", 200, fmt.Sprintf("%s's account has been deleted", customer.Email), nil)

	c.JSON(http.StatusOK, response)

}

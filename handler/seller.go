package handler

import (
	"errors"
	"fmt"
	"net/http"
	"olshop/auth"
	"olshop/product"
	"olshop/seller"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerSeller struct {
	service seller.ServiceInt
	auth    auth.Service
}

func NewSellerHandler(service seller.ServiceInt, auth auth.Service) *HandlerSeller {
	return &HandlerSeller{service: service, auth: auth}
}

func (h *HandlerSeller) RegisterSeller(c *gin.Context) {

	var input seller.InputSeller

	c.ShouldBindJSON(&input)
	if input.Password != input.ConfirmPassword {
		response := APIResponse("password and confirm password is different", http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	if len(input.Email) < 5 || len(input.Password) < 5 {
		response := APIResponse("your email or password's too short ", http.StatusForbidden, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	sellerSave := seller.Seller{}
	sellerSave.Name = input.Name
	sellerSave.Email = input.Email
	sellerSave.Password = input.Password

	seller, err := h.service.Register(sellerSave)

	if err != nil {
		respones := APIResponse("failed to create account", http.StatusBadRequest, fmt.Sprintf("%v", err.Error()), nil)
		c.JSON(http.StatusBadRequest, respones)
		return
	}

	token, err := h.auth.GenerateTokenSeller(seller.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("failed to generate token"))
		return

	}

	response := ResponseAPIToken("account successfully created", http.StatusOK, fmt.Sprintf("new seller suceessfully created with id %d", seller.ID), seller, token)
	c.JSON(http.StatusOK, response)
}

func (h *HandlerSeller) Login(c *gin.Context) {
	var input seller.InputLogin

	c.ShouldBindJSON(&input)

	seller, err := h.service.Login(input)
	if err != nil {
		respones := APIResponse("failed", http.StatusBadRequest, fmt.Sprintf("%v", err.Error()), nil)
		c.JSON(http.StatusBadRequest, respones)
		return
	}
	token, err := h.auth.GenerateTokenSeller(seller.ID)
	if err != nil {
		respones := APIResponse("failed to generate token", http.StatusBadRequest, fmt.Sprintf("%v", err.Error()), nil)
		c.JSON(http.StatusBadRequest, respones)
		return
	}

	response := ResponseAPIToken("success", http.StatusOK, fmt.Sprintf("%s's success login", seller.Email), seller, token)
	c.JSON(http.StatusOK, response)
}

func (h *HandlerSeller) InsertProduct(c *gin.Context) {
	productType := c.Request.FormValue("product_type")
	quantity, _ := strconv.Atoi(c.Request.FormValue("quantity"))
	price, _ := strconv.Atoi(c.Request.FormValue("price"))
	name := c.Request.FormValue("name")

	seller := c.MustGet("currentSeller").(seller.Seller)
	product := product.Product{}
	product.Quantity = quantity
	product.Price = int32(price)
	product.SellerID = seller.ID
	product.Name = name

	getProduct, err := h.service.InsertProduct(product, productType)
	if err != nil {
		respones := APIResponse("failed", http.StatusBadRequest, fmt.Sprintf("%v", err.Error()), nil)
		c.JSON(http.StatusBadRequest, respones)
		return
	}
	response := APIResponse("success", http.StatusOK, fmt.Sprintf("%s's success insert product", seller.Email), getProduct)
	c.JSON(http.StatusOK, response)

}

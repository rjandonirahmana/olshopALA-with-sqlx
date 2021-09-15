package handler

import (
	"fmt"
	"net/http"
	"olshop/customer"
	"olshop/product"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerProduct struct {
	service product.ServiceProductInt
}

func NewProductHandler(service product.ServiceProductInt) *HandlerProduct {
	return &HandlerProduct{service: service}
}

func (h *HandlerProduct) GetProductByCategory(c *gin.Context) {
	category := c.Request.FormValue("category")

	products, err := h.service.GetProductCategory(category)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("success", http.StatusOK, "products", products)
	c.JSON(http.StatusOK, response)

}

func (h *HandlerProduct) CreateShopCart(c *gin.Context) {

	customer := c.MustGet("currentCustomer").(customer.Customer)

	id, err := h.service.AddShoppingCart(customer.ID)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("success", http.StatusOK, fmt.Sprintf("Create cart for %s with shopcart id %d", customer.Email, id), nil)
	c.JSON(http.StatusOK, response)

}

func (h *HandlerProduct) InsertToShopCart(c *gin.Context) {

	idChart, err := strconv.Atoi((c.Request.FormValue("id")))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	productID, err := strconv.Atoi((c.Request.FormValue("product_id")))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	customer := c.MustGet("currentCustomer").(customer.Customer)

	product, err := h.service.InsertProductByCartID(customer.ID, productID, idChart)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusUnprocessableEntity, "failed", err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("success", http.StatusOK, fmt.Sprintf("successfully insert product %s to shopecart for %s", product.Name, customer.Email), product)
	c.JSON(http.StatusOK, response)

}

func (h *HandlerProduct) GetListProductShopCart(c *gin.Context) {

	cartID, err := strconv.Atoi((c.Request.FormValue("id")))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	customer := c.MustGet("currentCustomer").(customer.Customer)

	products, err := h.service.GetListInShopCart(cartID, customer.ID)

	if err != nil {
		response := APIResponse(err.Error(), http.StatusUnauthorized, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	if len(products) == 0 {
		response := APIResponse("success", http.StatusOK, fmt.Sprintf("%s' doesnt have any product in this cart,  cart_id : %d", customer.Email, cartID), products)
		c.JSON(http.StatusOK, response)
		return

	}
	response := APIResponse("success", http.StatusOK, fmt.Sprintf("list product shop cart with id customer : %d and cart_id : %d", customer.ID, cartID), products)
	c.JSON(http.StatusOK, response)

}

func (h *HandlerProduct) DeleteProductShopcart(c *gin.Context) {
	shopcartid, err := strconv.Atoi((c.Request.FormValue("shopcart_id")))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	productid, err := strconv.Atoi((c.Request.FormValue("product_id")))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	customer := c.MustGet("currentCustomer").(customer.Customer)

	productLeft, err := h.service.DeleteListOnshoppingCart(shopcartid, customer.ID, productid)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusConflict, "failed", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := APIResponse("success", http.StatusOK, fmt.Sprintf("list product left on shoppping cart with id customer : %d and shopcart_id : %d", customer.ID, shopcartid), productLeft)
	c.JSON(http.StatusOK, response)

}

func (h *HandlerProduct) GetAllCartCustomer(c *gin.Context) {
	customer := c.MustGet("currentCustomer").(customer.Customer)

	cartCustomer, err := h.service.GetShopCartCustomer(customer.ID)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusConflict, "failed", err)
		c.JSON(http.StatusInternalServerError, response)
		return

	}
	response := APIResponse("success", http.StatusOK, fmt.Sprintf("%s have %v shopcart", customer.Email, len(cartCustomer)), cartCustomer)
	c.JSON(http.StatusOK, response)

}

func (h *HandlerProduct) DecreaseQuantity(c *gin.Context) {
	cartid, err := strconv.Atoi((c.Request.FormValue("cart_id")))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	productID, err := strconv.Atoi((c.Request.FormValue("product_id")))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	customer := c.MustGet("currentCustomer").(customer.Customer)

	productLeft, err := h.service.DecreaseProductShopCart(customer.ID, productID, cartid)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusConflict, "failed", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := APIResponse("success", http.StatusOK, fmt.Sprintf("products left in shop cart with id customer : %d and cart_id : %d", customer.ID, cartid), productLeft)
	c.JSON(http.StatusOK, response)

}

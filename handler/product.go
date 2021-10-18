package handler

import (
	"net/http"
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	products, err := h.service.GetProductCategory(id)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("success", http.StatusOK, "products", products)
	c.JSON(http.StatusOK, response)

}

func (h *HandlerProduct) GetProductByID(c *gin.Context) {
	product_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := APIResponse("failed", http.StatusUnprocessableEntity, err.Error(), nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	product, err := h.service.GetProductByid(product_id)
	if err != nil {
		response := APIResponse("failed", http.StatusUnprocessableEntity, err.Error(), nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if product.ID == 0 {
		response := APIResponse("failed", http.StatusBadRequest, "product not found", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := APIResponse("success", http.StatusOK, "success get product", product)
	c.JSON(http.StatusOK, response)

}

package handler

import (
	"fmt"
	"net/http"
	"olshop/customer"
	"olshop/transaction"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerTransaction struct {
	service transaction.ServiceTransInt
}

func (h *HandlerTransaction) CreateTransaction(c *gin.Context) {

	quantity, err := strconv.Atoi(c.Request.FormValue("quantity"))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusBadRequest, "failed", err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	id, err := strconv.Atoi(c.Request.FormValue("product_id"))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusBadRequest, "failed", err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentCustomer := c.MustGet("currentCustomer").(customer.Customer)
	idCustomer := currentCustomer.ID

	trans, err := h.service.CreateTransaction(id, idCustomer, quantity)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusBadRequest, "failed", err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("success create transaction", 200, fmt.Sprintf("successfully create transaction id : %d", trans.ID), trans)

	c.JSON(http.StatusOK, response)

}

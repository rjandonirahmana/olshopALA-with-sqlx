package handler

// import (
// 	"fmt"
// 	"net/http"
// 	"olshop/customer"
// 	"olshop/transaction"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// type HandlerTransaction struct {
// 	service transaction.ServiceTransInt
// }

// func NewTransactionHandler(service transaction.ServiceTransInt) *HandlerTransaction {
// 	return &HandlerTransaction{service: service}
// }

// func (h *HandlerTransaction) CreateTransaction(c *gin.Context) {

// 	id, err := strconv.Atoi(c.Request.FormValue("cart_id"))
// 	if err != nil {
// 		response := APIResponse(err.Error(), http.StatusBadRequest, "failed", err)
// 		c.JSON(http.StatusUnprocessableEntity, response)
// 		return
// 	}

// 	currentCustomer := c.MustGet("currentCustomer").(customer.Customer)

// 	trans, err := h.service.CreateTransaction(currentCustomer, id)
// 	if err != nil {
// 		fmt.Println("error disinin2")
// 		response := APIResponse(err.Error(), http.StatusBadRequest, "failed", err)
// 		c.JSON(http.StatusUnprocessableEntity, response)
// 		return
// 	}

// 	response := APIResponse("success create transaction", 200, fmt.Sprintf("%s' successfully create transaction with id transaction: %d", currentCustomer.Email, trans.ID), trans)

// 	c.JSON(http.StatusOK, response)

// }

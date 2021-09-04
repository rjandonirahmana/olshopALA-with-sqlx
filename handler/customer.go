package handler

// type handlerCustomer struct {
// 	usecase usecase.CustomerInt
// }

// func NewHandlerCustomer(use usecase.CustomerInt) *handlerCustomer {
// 	return &handlerCustomer{usecase: use}
// }

// func (h *handlerCustomer) CreateCustomer(c *gin.Context) {
// 	var customer repo.Customer

// 	c.ShouldBindJSON(&customer)

// 	err := h.usecase.Register(customer)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, errors.New("gagal"))
// 		return
// 	}

// 	c.JSON(http.StatusOK, nil)

// }

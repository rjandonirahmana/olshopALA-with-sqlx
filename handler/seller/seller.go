package seller

import (
	"net/http"
	"olshop/auth"
	"olshop/handler"
	"olshop/seller"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HandlerSeller struct {
	service seller.Service
	auth    auth.Service
}

func NewHandlerSeller(s seller.Service, auth auth.Service) *HandlerSeller {
	return &HandlerSeller{service: s, auth: auth}
}

func (h *HandlerSeller) Register(c *gin.Context) {
	var input seller.InputSeller
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := handler.APIResponse("make sure to input whole field correctly", http.StatusUnprocessableEntity, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	validate := validator.New()
	err = validate.Struct(&input)
	if err != nil {
		response := handler.APIResponse("check again your input, password min 8character long and make sure your confirm password match ", http.StatusUnprocessableEntity, "failed", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	seller, err := h.service.Register(input)
	if err != nil {
		response := handler.APIResponse("error to register seller", http.StatusUnprocessableEntity, "failed", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.auth.GenerateTokenSeller(seller.ID)
	if err != nil {
		response := handler.APIResponse("error to generate token", http.StatusUnprocessableEntity, "failed", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	cookie := http.Cookie{
		Name:    "tokenseller",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 12),
	}

	http.SetCookie(c.Writer, &cookie)
	response := handler.ResponseAPIToken("new seller successfully created", http.StatusOK, "suceess", seller, token)
	c.JSON(http.StatusOK, response)

}

func (h *HandlerSeller) Login(c *gin.Context) {
	var input seller.InputLoginSeller

	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := handler.APIResponse("make sure to input whole field correctly", http.StatusUnprocessableEntity, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	validate := validator.New()
	err = validate.Struct(&input)
	if err != nil {
		response := handler.APIResponse("check again your input, password min 8character long and make sure your confirm password match ", http.StatusUnprocessableEntity, "failed", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	seller, err := h.service.LoginSeller(input)
	if err != nil {
		response := handler.APIResponse("error login seller", http.StatusUnprocessableEntity, "failed", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.auth.GenerateTokenSeller(seller.ID)
	if err != nil {
		response := handler.APIResponse("error to generate token", http.StatusUnprocessableEntity, "failed", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	cookie := http.Cookie{
		Name:    "tokenseller",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 12),
	}
	http.SetCookie(c.Writer, &cookie)
	response := handler.ResponseAPIToken("login success", http.StatusOK, "suceess", seller, token)
	c.JSON(http.StatusOK, response)
}

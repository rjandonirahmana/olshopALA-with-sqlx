package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"olshop/auth"
	"olshop/customer"
	"olshop/handler"
	"olshop/product"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/tj/assert"
)

func connectDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "root:12345@(localhost:3306)/dbsampingan?parseTime=true")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setupRouter() *gin.Engine {
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Error %s when opening DB", err)
		return nil
	}

	auth := auth.NewService()
	customerdb := customer.NewRepo(db)
	productdb := product.NewRepoProduct(db)
	customerserv := customer.NewCustomerService(customerdb)
	productServ := product.NewService(productdb)

	productHanlder := handler.NewProductHandler(productServ)
	customerHandler := handler.NewHandlerCustomer(customerserv, auth)

	gin.SetMode(gin.TestMode)
	c := gin.Default()
	c.GET("productCategory", productHanlder.GetProductByCategory)
	c.POST("/register", customerHandler.CreateCustomer)
	c.POST("/login", customerHandler.Login)
	c.PUT("/phone", customerHandler.UpdatePhoneCustomer)
	c.PUT("/avatar", customerHandler.UpdateAvatar)
	c.POST("/addcart", productHanlder.CreateShopCart)
	c.POST("/addshopcart", productHanlder.InsertToShopCart)
	c.GET("listshopcart", productHanlder.GetListProductShopCart)

	return c
}

func TestCreateCustomer(t *testing.T) {
	testCases := []struct {
		testName        string
		name            string
		email           string
		password        string
		confirmPassword string
		expectCode      int
		expectMsg       string
	}{
		{
			testName:        "success",
			name:            "joni2",
			email:           "SI KOCAK",
			password:        "jiji1",
			confirmPassword: "jiji1",
			expectCode:      http.StatusOK,
			expectMsg:       "account successfully created",
		},
		{
			testName:        "fail used email",
			name:            "joni3",
			email:           "ateng",
			password:        "jojo1",
			confirmPassword: "jojo1",
			expectCode:      http.StatusBadRequest,
			expectMsg:       "email has been used",
		},
		{
			testName:        "fail different password",
			name:            "joni4",
			email:           "jon@l.l",
			password:        "jojo",
			confirmPassword: "jiji",
			expectCode:      http.StatusUnprocessableEntity,
			expectMsg:       "password and confirm password is different",
		},
	}

	r := setupRouter()

	for _, testCase := range testCases {
		reqBody, _ := json.Marshal(map[string]string{
			"name":             testCase.name,
			"email":            testCase.email,
			"password":         testCase.password,
			"confirm_password": testCase.confirmPassword,
		})
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		fmt.Println(testCase.testName, res.Body.String())
		assert.Equal(t, testCase.expectCode, res.Code)
		// assert.True(t, strings.Contains(res.Body.String(), testCase.name))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

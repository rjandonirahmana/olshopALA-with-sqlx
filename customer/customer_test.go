package customer_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"olshop/auth"
	"olshop/customer"
	"olshop/handler"
	"olshop/product"
	"log"
	"net/http"
	"net/http/httptest"
	"olshop/transaction"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading environment variables")
		os.Exit(1)
	}
}

func connectDB() (*sqlx.DB, error) {
	LoadEnv()
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
		os.Getenv("DB_username"),
		os.Getenv("DB_password"),
		os.Getenv("DB_host"),
		os.Getenv("DB_port"),
		os.Getenv("DB_name"))
	db, err := sqlx.Connect("mysql", conn)
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
	transactiondb := transaction.NewTransactionRepo(db)

	customerserv := customer.NewCustomerService(customerdb)
	productServ := product.NewService(productdb)
	transactionServ := transaction.NewTransactionService(transactiondb, productdb)

	productHanlder := handler.NewProductHandler(productServ)
	customerHandler := handler.NewHandlerCustomer(customerserv, auth)
	transactionHandler := handler.NewTransactionHandler(transactionServ)

	gin.SetMode(gin.TestMode)
	c := gin.Default()

	c.GET("/productcategory", productHanlder.GetProductByCategory)
	c.POST("/register", customerHandler.CreateCustomer)
	c.POST("/login", customerHandler.Login)
	c.PUT("/phone", customerHandler.UpdatePhoneCustomer)
	c.PUT("/avatar", customerHandler.UpdateAvatar)
	c.PUT("password", customerHandler.UpdatePassword)
	c.DELETE("/account", customerHandler.DeleteAccount)
	c.POST("/addcart", productHanlder.CreateShopCart)

	c.POST("/insertshopcart", productHanlder.InsertToShopCart)
	c.GET("/listshopcart", productHanlder.GetListProductShopCart)
	c.GET("/shopcartcustomer", productHanlder.GetAllCartCustomer)
	c.PUT("/decreaseproduct", productHanlder.DecreaseQuantity)
	c.DELETE("/productshopcart", productHanlder.DeleteProductShopcart)
	c.POST("/transaction", transactionHandler.CreateTransaction)

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
			email:           "jon@k.k",
			password:        "jiji",
			confirmPassword: "jiji",
			expectCode:      http.StatusOK,
			expectMsg:       "account successfully created",
		},
		{
			testName:        "fail used email",
			name:            "joni3",
			email:           "jon@k.k",
			password:        "jojo",
			confirmPassword: "jojo",
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
		assert.True(t, strings.Contains(res.Body.String(), testCase.name))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestLogin(t *testing.T) {
	testCases := []struct {
		testName   string
		email      string
		password   string
		expectCode int
		expectMsg  string
	}{
		{
			testName:   "success",
			email:      "jon@k.k",
			password:   "jiji",
			expectCode: http.StatusOK,
			expectMsg:  "success",
		},
		{
			testName:   "fail email not found",
			email:      "jin@k.k",
			password:   "jojo",
			expectCode: http.StatusUnprocessableEntity,
			expectMsg:  "email not found",
		},
		{
			testName:   "fail wrong password",
			email:      "jon@k.k",
			password:   "jojo",
			expectCode: http.StatusUnprocessableEntity,
			expectMsg:  "different password",
		},
	}

	r := setupRouter()

	for _, testCase := range testCases {
		reqBody, _ := json.Marshal(map[string]string{
			"email":    testCase.email,
			"password": testCase.password,
		})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		fmt.Println(testCase.testName, res.Body.String())
		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.email))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestUpdatePhoneCustomer(t *testing.T) {
	// Testcases list
	testCases := []struct {
		testName   string
		email      string
		phone      string
		expectCode int
		expectMsg  string
	}{
		{
			testName:   "success",
			email:      "jon@k.k",
			phone:      `0812345`,
			expectCode: http.StatusOK,
			expectMsg:  "success",
		},
		{
			testName:   "fail email not found",
			email:      "jin@k.k",
			phone:      `0812345`,
			expectCode: http.StatusUnprocessableEntity,
			expectMsg:  "email not found",
		},
	}

	// setting handler
	auth := auth.NewService()
	db, _ := connectDB()
	r := customer.NewRepo(db)
	s := customer.NewCustomerService(r)
	h := handler.NewHandlerCustomer(s, auth)

	for _, testCase := range testCases {
		currentCustomer := customer.Customer{
			Email: testCase.email,
		}
		
		reqBody := fmt.Sprintf(`phone=%s`, testCase.phone)
		res := httptest.NewRecorder()
		
		req := httptest.NewRequest(http.MethodPut, "/phone", strings.NewReader(reqBody))
		
		c, r := gin.CreateTestContext(res)
		c.Request = req
		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
		c.Set("currentCustomer", currentCustomer)
		r.ServeHTTP(res, req)
		r.PUT("/phone", h.UpdatePhoneCustomer)
		h.UpdatePhoneCustomer(c)
		fmt.Println(res.Body.String())

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.email))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}

}

func TestUpdatePassword(t *testing.T) {
	// Testcases list
	testCases := []struct {
		testName   string
		id      int
		oldPassword      string
		newPassword string
		expectCode int
		expectMsg  string
	}{
		{
			testName:   "success",
			id:      1,
			oldPassword: "jojo",
			newPassword: "jiji",
			expectCode: http.StatusOK,
			expectMsg:  "success",
		},
		{
			testName:   "fail email not found",
			id:      2,
			oldPassword: "aaaa",
			newPassword: "uuu",
			expectCode: http.StatusUnprocessableEntity,
			expectMsg:  "please",
		},
	}

	// setting handler
	auth := auth.NewService()
	db, _ := connectDB()
	r := customer.NewRepo(db)
	s := customer.NewCustomerService(r)
	h := handler.NewHandlerCustomer(s, auth)

	for _, testCase := range testCases {
		currentCustomer := customer.Customer{
			ID: testCase.id,
		}
		
		reqBody := fmt.Sprintf(`password=%s&newpassword=%s`, testCase.oldPassword, testCase.newPassword)
		res := httptest.NewRecorder()
		
		req := httptest.NewRequest(http.MethodPut, "/password", strings.NewReader(reqBody))
		
		c, r := gin.CreateTestContext(res)
		c.Request = req
		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
		c.Set("currentCustomer", currentCustomer)
		r.ServeHTTP(res, req)
		r.PUT("/password", h.UpdatePassword)
		h.UpdatePassword(c)
		fmt.Println(res.Body.String())

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestDeleteAccount(t *testing.T) {
	// Testcases list
	testCases := []struct {
		testName   string
		id      int
		password      string
		expectCode int
		expectMsg  string
	}{
		{
			testName:   "success",
			id:      1,
			password:      "jojo",
			expectCode: http.StatusOK,
			expectMsg:  "success",
		},
		{
			testName:   "fail email not found",
			id:      2,
			password:      "aaa",
			expectCode: http.StatusUnprocessableEntity,
			expectMsg:  "cant delete",
		},
	}

	// setting handler
	auth := auth.NewService()
	db, _ := connectDB()
	r := customer.NewRepo(db)
	s := customer.NewCustomerService(r)
	h := handler.NewHandlerCustomer(s, auth)

	for _, testCase := range testCases {
		currentCustomer := customer.Customer{
			ID: testCase.id,
		}
		
		reqBody := fmt.Sprintf(`password=%s`, testCase.password)
		res := httptest.NewRecorder()
		
		req := httptest.NewRequest(http.MethodDelete, "/account", strings.NewReader(reqBody))
		
		c, r := gin.CreateTestContext(res)
		c.Request = req
		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
		c.Set("currentCustomer", currentCustomer)
		r.ServeHTTP(res, req)
		r.DELETE("/account", h.DeleteAccount)
		h.DeleteAccount(c)
		fmt.Println(res.Body.String())

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}
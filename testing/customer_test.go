package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"olshop/auth"
	"olshop/customer"
	"olshop/handler"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	response handler.Response
)

func setupRouter() *gin.Engine {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading environment variables")
		os.Exit(1)
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
		os.Getenv("DB_username"),
		os.Getenv("DB_password"),
		os.Getenv("DB_host"),
		os.Getenv("DB_port"),
		os.Getenv("DB_name"))
	db, err := sqlx.Connect("mysql", conn)
	if err != nil {
		panic(err)
	}

	auth := auth.NewService("coba", "cobalagi")
	customerdb := customer.NewRepo(db)
	customerserv := customer.NewCustomerService(customerdb)
	customerHandler := handler.NewHandlerCustomer(customerserv, auth)

	gin.SetMode(gin.TestMode)
	c := gin.Default()
	api := c.Group("/api/v1")

	api.POST("/register", customerHandler.CreateCustomer)
	api.POST("/login", customerHandler.Login)
	api.PUT("/phone", customerHandler.UpdatePhoneCustomer)
	api.PUT("/avatar", customerHandler.UpdateAvatar)
	api.PUT("password", customerHandler.UpdatePassword)
	api.DELETE("/account", customerHandler.DeleteAccount)

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
			testName:        "test1",
			name:            "joni2",
			email:           "doniaaaa12@gmail.com",
			password:        "jijijiji",
			confirmPassword: "jijijiji",
			expectCode:      http.StatusOK,
			expectMsg:       "new customer successfully created",
		},
		{
			testName:        "test2",
			name:            "joni3",
			email:           "baruyaaa@gmail.com",
			password:        "jojojojo",
			confirmPassword: "jojojojo",
			expectCode:      http.StatusForbidden,
			expectMsg:       "email has been used",
		},
		{
			testName:        "test3",
			name:            "joni4",
			email:           "jon@l.lala",
			password:        "jojo",
			confirmPassword: "jiji",
			expectCode:      http.StatusForbidden,
			expectMsg:       "make sure to input whole field correctly",
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
		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		fmt.Printf("====================================================================%s===============================================", testCase.testName)
		fmt.Println(testCase.testName, res.Body.String())

		resBody := response
		err := json.Unmarshal(res.Body.Bytes(), &resBody)
		assert.NoError(t, err)

		fmt.Println(resBody.Meta.Status)
		assert.True(t, strings.Contains(resBody.Meta.Status, testCase.expectMsg))
		assert.Equal(t, testCase.expectCode, resBody.Meta.Code)
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
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
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
			testName:   "test1",
			email:      "jon@k.k",
			phone:      `0812345`,
			expectCode: http.StatusOK,
			expectMsg:  "success",
		},
		{
			testName:   "test2",
			email:      "jin@k.k",
			phone:      `0812345`,
			expectCode: http.StatusUnprocessableEntity,
			expectMsg:  "email not found",
		},
	}

	for _, testCase := range testCases {
		currentCustomer := customer.Customer{
			Email: testCase.email,
		}

		reqBody := fmt.Sprintf(`phone=%s`, testCase.phone)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/phone", strings.NewReader(reqBody))
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
		c.Set("currentCustomer", currentCustomer)
		user := c.MustGet("currentCustomer").(customer.Customer)
		c.Request = req
		fmt.Println(user)

		fmt.Println(res)
		fmt.Println(res.Body.String())

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.email))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}

}

// func TestUpdatePassword(t *testing.T) {
// 	// Testcases list
// 	testCases := []struct {
// 		testName    string
// 		id          int
// 		oldPassword string
// 		newPassword string
// 		expectCode  int
// 		expectMsg   string
// 	}{
// 		{
// 			testName:    "success",
// 			id:          1,
// 			oldPassword: "jojo",
// 			newPassword: "jiji",
// 			expectCode:  http.StatusOK,
// 			expectMsg:   "success",
// 		},
// 		{
// 			testName:    "fail email not found",
// 			id:          2,
// 			oldPassword: "aaaa",
// 			newPassword: "uuu",
// 			expectCode:  http.StatusUnprocessableEntity,
// 			expectMsg:   "please",
// 		},
// 	}

// 	// setting handler
// 	auth := auth.NewService()
// 	db, _ := connectDB()
// 	r := customer.NewRepo(db)
// 	s := customer.NewCustomerService(r)
// 	h := handler.NewHandlerCustomer(s, auth)

// 	for _, testCase := range testCases {
// 		currentCustomer := customer.Customer{
// 			ID: testCase.id,
// 		}

// 		reqBody := fmt.Sprintf(`password=%s&newpassword=%s`, testCase.oldPassword, testCase.newPassword)
// 		res := httptest.NewRecorder()

// 		req := httptest.NewRequest(http.MethodPut, "/password", strings.NewReader(reqBody))

// 		c, r := gin.CreateTestContext(res)
// 		c.Request = req
// 		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
// 		c.Set("currentCustomer", currentCustomer)
// 		r.ServeHTTP(res, req)
// 		r.PUT("/password", h.UpdatePassword)
// 		h.UpdatePassword(c)
// 		fmt.Println(res.Body.String())

// 		assert.Equal(t, testCase.expectCode, res.Code)
// 		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
// 	}
// }

// func TestDeleteAccount(t *testing.T) {
// 	// Testcases list
// 	testCases := []struct {
// 		testName   string
// 		id         int
// 		password   string
// 		expectCode int
// 		expectMsg  string
// 	}{
// 		{
// 			testName:   "success",
// 			id:         1,
// 			password:   "jojo",
// 			expectCode: http.StatusOK,
// 			expectMsg:  "success",
// 		},
// 		{
// 			testName:   "fail email not found",
// 			id:         2,
// 			password:   "aaa",
// 			expectCode: http.StatusUnprocessableEntity,
// 			expectMsg:  "cant delete",
// 		},
// 	}

// 	// setting handler
// 	auth := auth.NewService()
// 	db, _ := connectDB()
// 	r := customer.NewRepo(db)
// 	s := customer.NewCustomerService(r)
// 	h := handler.NewHandlerCustomer(s, auth)

// 	for _, testCase := range testCases {
// 		currentCustomer := customer.Customer{
// 			ID: testCase.id,
// 		}

// 		reqBody := fmt.Sprintf(`password=%s`, testCase.password)
// 		res := httptest.NewRecorder()

// 		req := httptest.NewRequest(http.MethodDelete, "/account", strings.NewReader(reqBody))

// 		c, r := gin.CreateTestContext(res)
// 		c.Request = req
// 		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
// 		c.Set("currentCustomer", currentCustomer)
// 		r.ServeHTTP(res, req)
// 		r.DELETE("/account", h.DeleteAccount)
// 		h.DeleteAccount(c)
// 		fmt.Println(res.Body.String())

// 		assert.Equal(t, testCase.expectCode, res.Code)
// 		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
// 	}
// }

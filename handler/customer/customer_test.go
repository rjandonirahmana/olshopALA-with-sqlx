package customer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	"olshop/auth"
	"olshop/customer"
	"olshop/handler"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupdb() *sqlx.DB {
	err := godotenv.Load("../../.env")
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

	return db
}

var (
	response        handler.Response
	authorization   = auth.NewService("coba", "cobalagi")
	db              = setupdb()
	customerdb      = customer.NewRepo(db)
	customerserv    = customer.NewCustomerService(customerdb)
	customerHandler = NewHandlerCustomer(customerserv, authorization)
	middlerware     = handler.NewMiddleWare()
)

func TestCreateCustomer(t *testing.T) {
	testCases := []struct {
		testName        string
		name            string
		email           string
		password        interface{}
		confirmPassword string
		expectCode      int
		expectMsg       string
	}{
		{
			testName:        "test1",
			name:            "joni2",
			email:           "doniaaaa12@gmail.com",
			password:        "123456",
			confirmPassword: "123456",
			expectCode:      http.StatusForbidden,
			expectMsg:       "your password need to be 8 length and make sure your confirm password match your password",
		},
		{
			testName:        "test2",
			name:            "joni3",
			email:           "baruyaaa@gmail.com",
			password:        123455,
			confirmPassword: "12345",
			expectCode:      http.StatusUnprocessableEntity,
			expectMsg:       "make sure to input whole field correctly",
		},
		{
			testName:        "test3",
			name:            "joni4",
			email:           "jon@gmail.com",
			password:        "jojojojo",
			confirmPassword: "jojojojo",
			expectCode:      http.StatusOK,
			expectMsg:       "new customer successfully created with id",
		},
	}

	for _, testCase := range testCases {
		reqBody, _ := json.Marshal(map[string]interface{}{
			"name":             testCase.name,
			"email":            testCase.email,
			"password":         testCase.password,
			"confirm_password": testCase.confirmPassword,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = req
		customerHandler.CreateCustomer(c)
		fmt.Printf("====================================================================%s===============================================", testCase.testName)
		fmt.Println(res)

		resBody := response
		err := json.Unmarshal(res.Body.Bytes(), &resBody)
		assert.NoError(t, err)

		fmt.Println(resBody.Meta.Status)
		assert.True(t, strings.Contains(resBody.Meta.Message, testCase.expectMsg))
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

	for _, testCase := range testCases {
		reqBody, _ := json.Marshal(map[string]string{
			"email":    testCase.email,
			"password": testCase.password,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = req
		customerHandler.Login(c)

		fmt.Println(res)
		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.email))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

// func TestUpdatePhoneCustomer(t *testing.T) {
// 	// Testcases list
// 	testCases := []struct {
// 		testName   string
// 		email      string
// 		phone      string
// 		expectCode int
// 		expectMsg  string
// 	}{
// 		{
// 			testName:   "test1",
// 			email:      "jon@k.k",
// 			phone:      `0812345`,
// 			expectCode: http.StatusOK,
// 			expectMsg:  "success",
// 		},
// 		{
// 			testName:   "test2",
// 			email:      "jin@k.k",
// 			phone:      `0812345`,
// 			expectCode: http.StatusUnprocessableEntity,
// 			expectMsg:  "email not found",
// 		},
// 	}

// 	for _, testCase := range testCases {

// 		fmt.Println("3")
// 		res := httptest.NewRecorder()
// 		fmt.Println("4")

// 		ctx, _ := gin.CreateTestContext(res)
// 		ctx.Header("Content-Type", "application/x-www-form-urlencoded")
// 		middlerware.AuthMiddleWareCustomer(authorization, customerserv)
// 		customerHandler.UpdatePhoneCustomer(ctx)

// 		fmt.Println(res)
// 		fmt.Println(testCase.testName)
// 		assert.NoError(t, nil)
// 		// fmt.Println(res.Body.String())

// 		// assert.Equal(t, testCase.expectCode, res.Code)
// 		// assert.True(t, strings.Contains(res.Body.String(), testCase.email))
// 		// assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
// 	}

// }

func TestUpdatePassword(t *testing.T) {
	// Testcases list
	testCases := []struct {
		testName    string
		id          int
		oldPassword string
		newPassword string
		expectCode  int
		expectMsg   string
	}{
		{
			testName:    "test1",
			id:          1,
			oldPassword: "jij",
			newPassword: "123456",
			expectCode:  http.StatusOK,
			expectMsg:   "success",
		},
		{
			testName:    "test2",
			id:          10,
			oldPassword: "uuu",
			newPassword: "123456",
			expectCode:  http.StatusUnprocessableEntity,
			expectMsg:   "please",
		},
	}

	for _, testCase := range testCases {

		customer, err := customerdb.GetCustomerByID(testCase.id)
		assert.NoError(t, err)

		f := make(url.Values)
		f.Set("password", testCase.oldPassword)
		f.Set("newpassword", testCase.newPassword)
		req := httptest.NewRequest(http.MethodPut, "/books/newauthor", strings.NewReader(f.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = req
		c.Set("currentCustomer", customer)
		customerHandler.UpdatePassword(c)

		fmt.Println(res)

		// assert.Equal(t, testCase.expectCode, res.Code)
		// assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestDeleteAccount(t *testing.T) {

	testCases := []struct {
		testName   string
		id         int
		password   string
		expectCode int
		expectMsg  string
	}{
		{
			testName:   "success",
			id:         1,
			password:   "jojo",
			expectCode: http.StatusOK,
			expectMsg:  "success",
		},
		{
			testName:   "fail email not found",
			id:         5,
			password:   "aaa",
			expectCode: http.StatusUnprocessableEntity,
			expectMsg:  "cant delete",
		},
	}

	for _, testCase := range testCases {
		customer, err := customerserv.GetCustomerByID(testCase.id)
		assert.NoError(t, err)

		f := make(url.Values)
		f.Add("password", testCase.password)
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/account", strings.NewReader(f.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		c, _ := gin.CreateTestContext(res)
		c.Request = req
		c.Set("currentCustomer", customer)
		customerHandler.DeleteAccount(c)
		fmt.Println(res)

		// assert.Equal(t, testCase.expectCode, res.Code)
		// assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

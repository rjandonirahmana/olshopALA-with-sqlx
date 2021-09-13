package customer_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"graphql/customer"
	"graphql/handler"
	"graphql/product"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
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

	customerdb := customer.NewRepo(db)
	productdb := product.NewRepoProduct(db)
	customerserv := customer.NewCustomerService(customerdb)
	productServ := product.NewService(productdb)

	productHanlder := handler.NewProductHandler(productServ)
	customerHandler := handler.NewHandlerCustomer(customerserv, nil)

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
			expectCode:      http.StatusBadRequest,
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
			expectCode: http.StatusUnprocessableEntity,
			expectMsg:  "success login",
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
// 			testName:   "success",
// 			email:      "jon@k.k",
// 			phone:      "0812345",
// 			expectCode: http.StatusUnprocessableEntity,
// 			expectMsg:  "successfully udpate",
// 		},
// 		{
// 			testName:   "fail email not found",
// 			email:      "jin@k.k",
// 			phone:      "0812345",
// 			expectCode: http.StatusUnprocessableEntity,
// 			expectMsg:  "email not found",
// 		},
// 	}

// 	// r := setupRouter()
// 	// setting handler
// 	db, _ := connectDB()
// 	r := customer.NewRepo(db)
// 	s := customer.NewCustomerService(r)
// 	h := handler.NewHandlerCustomer(s, nil)

// 	for _, testCase := range testCases {
// 		currentCustomer := customer.Customer{
// 			Email: testCase.email,
// 		}
		
// 		reqBody := url.Values{}
// 		reqBody.Set("phone", testCase.phone)

// 		res := httptest.NewRecorder()
		
// 		req := httptest.NewRequest(http.MethodPut, "/phone", strings.NewReader(reqBody.Encode()))
// 		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		
// 		c, r := gin.CreateTestContext(res)
// 		r.ServeHTTP(res, req)
// 		r.PUT("/phone", h.UpdatePhoneCustomer)
// 		c.Set("currentCustomer", currentCustomer)
// 		h.UpdatePhoneCustomer(c)

// 		assert.Equal(t, testCase.expectCode, res.Code)
// 		assert.True(t, strings.Contains(res.Body.String(), testCase.email))
// 		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
// 	}

// }

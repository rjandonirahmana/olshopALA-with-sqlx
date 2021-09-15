package product_test

// Insert dummy data first

import (
	"fmt"
	"olshop/customer"
	"olshop/handler"
	"olshop/product"
	"log"
	"net/http"
	"net/http/httptest"
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

func TestGetProductCategory(t *testing.T) {
	// Testcases list
	testCases := []struct{
		testName string
		category string
		expectCode int
		expectMsg string
	}{
		{
			testName: "success",
			category: "book",
			expectCode: http.StatusOK,
			expectMsg: "success",
		},
	}

	// setting handler
	db, _ := connectDB()
	r := product.NewRepoProduct(db)
	s := product.NewService(r)
	h := handler.NewProductHandler(s)

	for _, testCase := range testCases {
		reqBody := fmt.Sprintf(`category=%s`, testCase.category)
		req := httptest.NewRequest(http.MethodGet, "/productCategory", strings.NewReader(reqBody))
		res := httptest.NewRecorder()
		
		c, r := gin.CreateTestContext(res)
		c.Request = req
		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
		r.ServeHTTP(res, req)

		h.GetProductByCategory(c)
		r.GET("/productCategory", h.GetProductByCategory)
		h.GetProductByCategory(c)

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.category))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestAddShoppingCart(t *testing.T) {
	// Testcases list
	testCases := []struct{
		testName string
		customerId int
		id string
		expectCode int
		expectMsg string
	}{
		{
			testName: "success",
			customerId: 1,
			id: "1",
			expectCode: http.StatusOK,
			expectMsg: "success",
		},
	}

	// setting handler
	db, _ := connectDB()
	r := product.NewRepoProduct(db)
	s := product.NewService(r)
	h := handler.NewProductHandler(s)

	for _, testCase := range testCases {
		currentCustomer := customer.Customer{
			ID: testCase.customerId,
		}
		req := httptest.NewRequest(http.MethodPost, "/addcart", nil)
		res := httptest.NewRecorder()
		c, r := gin.CreateTestContext(res)
		c.Set("currentCustomer", currentCustomer)
		r.ServeHTTP(res, req)
		r.POST("/addcart", h.CreateShopCart)
		h.CreateShopCart(c)

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.id))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestInsertProductByCartId(t *testing.T) {
	// Testcases list
	testCases := []struct{
		testName string
		customerId int
		id string
		productId string
		expectCode int
		expectMsg string
	}{
		{
			testName: "success",
			customerId: 1,
			id: "1",
			productId: "1",
			expectCode: http.StatusOK,
			expectMsg: "success",
		},
	}

	// setting handler
	db, _ := connectDB()
	r := product.NewRepoProduct(db)
	s := product.NewService(r)
	h := handler.NewProductHandler(s)

	for _, testCase := range testCases {
		currentCustomer := customer.Customer{
			ID: testCase.customerId,
		}
		reqBody := fmt.Sprintf(`id=%s&product_id=%s`, testCase.id, testCase.productId)
		req := httptest.NewRequest(http.MethodPost, "/addshopcart", strings.NewReader(reqBody))
		res := httptest.NewRecorder()
		c, r := gin.CreateTestContext(res)
		c.Request = req
		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
		c.Set("currentCustomer", currentCustomer)
		r.ServeHTTP(res, req)
		r.POST("/addshopcart", h.InsertToShopCart)
		h.InsertToShopCart(c)

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.productId))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestGetListShopCart(t *testing.T) {
	// Testcases list
	testCases := []struct{
		testName string
		customerId int
		id string
		productId string
		expectCode int
		expectMsg string
	}{
		{
			testName: "success",
			customerId: 1,
			id: "1",
			productId: "1",
			expectCode: http.StatusOK,
			expectMsg: "success",
		},
	}

	// setting handler
	db, _ := connectDB()
	r := product.NewRepoProduct(db)
	s := product.NewService(r)
	h := handler.NewProductHandler(s)

	for _, testCase := range testCases {
		currentCustomer := customer.Customer{
			ID: testCase.customerId,
		}
		reqBody := fmt.Sprintf(`id=%s`, testCase.id)
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/listshopcart", strings.NewReader(reqBody))
		c, r := gin.CreateTestContext(res)
		c.Request = req
		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
		c.Set("currentCustomer", currentCustomer)
		r.ServeHTTP(res, req)
		r.GET("/listshopcart", h.GetListProductShopCart)
		h.GetListProductShopCart(c)

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.productId))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestGetAllCart(t *testing.T) {
	// Testcases list
	testCases := []struct{
		testName string
		customerId int
		lenCart string
		expectCode int
		expectMsg string
	}{
		{
			testName: "success",
			customerId: 1,
			lenCart: "1",
			expectCode: http.StatusOK,
			expectMsg: "have",
		},
	}

	// setting handler
	db, _ := connectDB()
	r := product.NewRepoProduct(db)
	s := product.NewService(r)
	h := handler.NewProductHandler(s)

	for _, testCase := range testCases {
		currentCustomer := customer.Customer{
			ID: testCase.customerId,
		}
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/shopcartcustomer", nil)
		c, r := gin.CreateTestContext(res)
		c.Request = req
		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
		c.Set("currentCustomer", currentCustomer)
		r.ServeHTTP(res, req)
		r.GET("/shopcartcustomer", h.GetAllCartCustomer)
		h.GetAllCartCustomer(c)

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.lenCart))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}
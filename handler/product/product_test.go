package product

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"olshop/product"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var (
	db, _ = connectDB()
	r     = product.NewRepoProduct(db)
	s     = product.NewService(r)
	h     = NewProductHandler(s)
)

func LoadEnv() {
	err := godotenv.Load("../../.env")
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
	testCases := []struct {
		testName   string
		category   string
		expectCode int
		expectMsg  string
	}{
		{
			testName:   "success",
			category:   "6",
			expectCode: http.StatusOK,
			expectMsg:  "success",
		},
	}

	// setting handler

	for _, testCase := range testCases {

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/category", nil)
		req.Header.Set("Content-Type", "application/json")
		q := req.URL.Query()
		q.Add("id", testCase.category)
		req.URL.RawQuery = q.Encode()
		c.Request = req

		h.GetProductByCategory(c)

		fmt.Println(res)

	}
}

func TestGetProductID(t *testing.T) {
	testCases := []struct {
		name         string
		request      string
		expectedCode int
	}{
		{
			name:         "test1",
			request:      "2",
			expectedCode: 200,
		},
	}

	for _, test := range testCases {
		f := make(url.Values)
		f.Set("id", test.request)
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/product/", strings.NewReader(f.Encode()))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		h.GetProductByID(c)

		fmt.Println(res)

	}
}

// func TestAddShoppingCart(t *testing.T) {
// 	// Testcases list
// 	testCases := []struct {
// 		testName   string
// 		customerId int
// 		id         string
// 		expectCode int
// 		expectMsg  string
// 	}{
// 		{
// 			testName:   "success",
// 			customerId: 1,
// 			id:         "1",
// 			expectCode: http.StatusOK,
// 			expectMsg:  "success",
// 		},
// 	}

// 	// setting handler
// 	db, _ := connectDB()
// 	r := product.NewRepoProduct(db)
// 	s := product.NewService(r)
// 	h := handler.NewProductHandler(s)

// 	for _, testCase := range testCases {
// 		currentCustomer := customer.Customer{
// 			ID: testCase.customerId,
// 		}
// 		req := httptest.NewRequest(http.MethodPost, "/addcart", nil)
// 		res := httptest.NewRecorder()
// 		c, r := gin.CreateTestContext(res)
// 		c.Set("currentCustomer", currentCustomer)
// 		r.ServeHTTP(res, req)
// 		r.POST("/addcart", h.CreateShopCart)
// 		h.CreateShopCart(c)

// 		assert.Equal(t, testCase.expectCode, res.Code)
// 		assert.True(t, strings.Contains(res.Body.String(), testCase.id))
// 		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
// 	}
// }

// // func TestInsertProductByCartId(t *testing.T) {
// // 	// Testcases list
// // 	testCases := []struct {
// // 		testName   string
// // 		customerId int
// // 		id         string
// // 		productId  string
// // 		expectCode int
// // 		expectMsg  string
// // 	}{
// // 		{
// // 			testName:   "success",
// // 			customerId: 1,
// // 			id:         "1",
// // 			productId:  "1",
// // 			expectCode: http.StatusOK,
// // 			expectMsg:  "success",
// // 		},
// // 	}

// // 	// setting handler
// // 	db, _ := connectDB()
// // 	r := product.NewRepoProduct(db)
// // 	s := product.NewService(r)
// // 	h := handler.NewProductHandler(s)

// // 	for _, testCase := range testCases {
// // 		currentCustomer := customer.Customer{
// // 			ID: testCase.customerId,
// // 		}
// // 		reqBody := fmt.Sprintf(`id=%s&product_id=%s`, testCase.id, testCase.productId)
// // 		req := httptest.NewRequest(http.MethodPost, "/addshopcart", strings.NewReader(reqBody))
// // 		res := httptest.NewRecorder()
// // 		c, r := gin.CreateTestContext(res)
// // 		c.Request = req
// // 		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
// // 		c.Set("currentCustomer", currentCustomer)
// // 		r.ServeHTTP(res, req)
// // 		r.POST("/addshopcart", h.InsertToShopCart)
// // 		h.InsertToShopCart(c)

// // 		assert.Equal(t, testCase.expectCode, res.Code)
// // 		assert.True(t, strings.Contains(res.Body.String(), testCase.productId))
// // 		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
// // 	}
// // }

// func TestGetListShopCart(t *testing.T) {
// 	// Testcases list
// 	testCases := []struct {
// 		testName   string
// 		customerId int
// 		id         string
// 		productId  string
// 		expectCode int
// 		expectMsg  string
// 	}{
// 		{
// 			testName:   "success",
// 			customerId: 1,
// 			id:         "1",
// 			productId:  "1",
// 			expectCode: http.StatusOK,
// 			expectMsg:  "success",
// 		},
// 	}

// 	// setting handler
// 	db, _ := connectDB()
// 	r := product.NewRepoProduct(db)
// 	s := product.NewService(r)
// 	h := handler.NewProductHandler(s)

// 	for _, testCase := range testCases {
// 		currentCustomer := customer.Customer{
// 			ID: testCase.customerId,
// 		}
// 		reqBody := fmt.Sprintf(`id=%s`, testCase.id)
// 		res := httptest.NewRecorder()
// 		req := httptest.NewRequest(http.MethodGet, "/listshopcart", strings.NewReader(reqBody))
// 		c, r := gin.CreateTestContext(res)
// 		c.Request = req
// 		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
// 		c.Set("currentCustomer", currentCustomer)
// 		r.ServeHTTP(res, req)
// 		r.GET("/listshopcart", h.GetListProductShopCart)
// 		h.GetListProductShopCart(c)

// 		assert.Equal(t, testCase.expectCode, res.Code)
// 		assert.True(t, strings.Contains(res.Body.String(), testCase.productId))
// 		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
// 	}
// }
// }

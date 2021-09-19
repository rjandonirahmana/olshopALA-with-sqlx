package product_test

// Insert dummy data first

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"olshop/customer"
	"olshop/handler"
	"olshop/product"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func connectDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "root:12345@(localhost:3306)/olshopALA?parseTime=true")
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
			category:   "gadget",
			expectCode: http.StatusOK,
			expectMsg:  "success",
		},
	}

	// setting handler
	db, _ := connectDB()
	r := product.NewRepoProduct(db)
	s := product.NewService(r)
	h := handler.NewProductHandler(s)

	for _, testCase := range testCases {
		reqBody := fmt.Sprintf(`category=%s`, testCase.category)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/productcategory", strings.NewReader(reqBody))
		res := httptest.NewRecorder()

		c, r := gin.CreateTestContext(res)
		c.Request = req
		c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)
		r.Group("/api/v1")
		h.GetProductByCategory(c)
		r.GET("/productcategory", h.GetProductByCategory)
		r.ServeHTTP(res, req)

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.category))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestAddShoppingCart(t *testing.T) {
	// Testcases list
	testCases := []struct {
		testName   string
		customerId int
		id         string
		expectCode int
		expectMsg  string
	}{
		{
			testName:   "success",
			customerId: 1,
			id:         "1",
			expectCode: http.StatusOK,
			expectMsg:  "success",
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
		req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
		res := httptest.NewRecorder()
		c, r := gin.CreateTestContext(res)

		r.POST("/addcart", h.CreateShopCart)
		c.Set("currentCustomer", currentCustomer)
		h.CreateShopCart(c)
		r.ServeHTTP(res, req)
		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.id))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestInsertProductByCartId(t *testing.T) {
	// Testcases list
	testCases := []struct {
		testName   string
		customerId int
		id         string
		productId  string
		expectCode int
		expectMsg  string
	}{
		{
			testName:   "success",
			customerId: 1,
			id:         "1",
			productId:  "1",
			expectCode: http.StatusOK,
			expectMsg:  "success",
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
	testCases := []struct {
		testName   string
		customerId int
		id         string
		productId  string
		expectCode int
		expectMsg  string
	}{
		{
			testName:   "success",
			customerId: 1,
			id:         "1",
			productId:  "1",
			expectCode: http.StatusOK,
			expectMsg:  "success",
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

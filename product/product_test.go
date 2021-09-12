package product_test

// Insert dummy data first

import (
	"fmt"
	"graphql/customer"
	"graphql/handler"
	"graphql/product"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables")
		os.Exit(1)
	}
}

// func TestMain(m *testing.M) {
// 	setup()
// 	gin.SetMode(gin.TestMode)
// 	os.Exit(m.Run())
// }

func connectDB(dbName string) (*sqlx.DB, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_username"),
		os.Getenv("DB_password"),
		os.Getenv("DB_host"),
		os.Getenv("DB_port"),
		dbName)
	db, err := sqlx.Connect("mysql", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// func setup() {
// 	// Create database connection
// 	LoadEnv()
// 	db, err := connectDB("")
// 	if err != nil {
// 		log.Printf("Error %s when opening DB", err)
// 		return
// 	}

// 	dbname := os.Getenv("DB_name")

// 	// Cleaning database before testing
// 	db.Exec("DROP DATABASE IF EXISTS " + dbname)

// 	// Create test database
// 	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
// 	if err != nil {
// 		log.Printf("Error %s when creating DB\n", err)
// 		return
// 	}
// 	db.Close()

// 	// Create database connection
// 	db, err = connectDB(dbname)
// 	if err != nil {
// 		log.Printf("Error %s when opening DB", err)
// 		return
// 	}

// 	// Create table
// 	db.MustExec(`CREATE TABLE products(
// 							name varchar(30),
// 							id int,
// 							price int32,
// 							category_id int,
// 							product_images )`)
// 	// if err != nil {
// 	// 	log.Printf("Error %s when creating table", err)
// 	// 	return
// 	// }

// 	// Insert dummy data
// 	customer := customer.Customer{}
// 	customer.ID = 1
// 	customer.Name = "joni"
// 	customer.Email = "jon@j.j"
// 	customer.Password = "jojo"
// 	customer.CreatedAt = time.Now()
// 	customer.UpdatedAt = time.Now()

// 	_, err = db.Exec(`INSERT INTO customers
// 				(id, name, email, password, created_at, updated_at)
// 				VALUES (?, ?, ?, ?, ?, ?)`,
// 		customer.ID, customer.Name, customer.Email,
// 		customer.Password, customer.CreatedAt, customer.UpdatedAt)
// 	if err != nil {
// 		log.Printf("Error %s when inserting dummy data", err)
// 	}
// }

func TestGetProductCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	LoadEnv()
	// Create database connection
	dbname := os.Getenv("DB_name")
	db, err := connectDB(dbname)
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}

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
	r := product.NewRepoProduct(db)
	s := product.NewService(r)
	h := handler.NewProductHandler(s)

	for _, testCase := range testCases {
		reqBody := url.Values{}
		reqBody.Set("category", testCase.category)
		g := gin.New()
		req := httptest.NewRequest(http.MethodGet, "productCategory", strings.NewReader(reqBody.Encode()))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		g.ServeHTTP(res, req)
		c, _ := gin.CreateTestContext(res)
		h.GetProductByCategory(c)

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.category))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestAddShoppingCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	LoadEnv()

	// Create database connection
	dbname := os.Getenv("DB_name")
	db, err := connectDB(dbname)
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}

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
	r := product.NewRepoProduct(db)
	s := product.NewService(r)
	h := handler.NewProductHandler(s)

	for _, testCase := range testCases {
		currentCustomer := customer.Customer{
			ID: testCase.customerId,
		}
		g := gin.New()
		req := httptest.NewRequest(http.MethodPost, "/addcart", nil)
		res := httptest.NewRecorder()
		g.ServeHTTP(res, req)
		c, _ := gin.CreateTestContext(res)
		c.Set("currentCustomer", currentCustomer)
		h.CreateShopCart(c)

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.id))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestInsertProductByCartId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	LoadEnv()
	
	// Create database connection
	dbname := os.Getenv("DB_name")
	db, err := connectDB(dbname)
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}

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
	r := product.NewRepoProduct(db)
	s := product.NewService(r)
	h := handler.NewProductHandler(s)

	for _, testCase := range testCases {
		currentCustomer := customer.Customer{
			ID: testCase.customerId,
		}
		reqBody := url.Values{}
		reqBody.Set("id", testCase.id)
		reqBody.Set("product_id", testCase.productId)
		g := gin.New()
		req := httptest.NewRequest(http.MethodPost, "/addshopcart", strings.NewReader(reqBody.Encode()))
		res := httptest.NewRecorder()
		g.ServeHTTP(res, req)
		c, _ := gin.CreateTestContext(res)
		c.Set("currentCustomer", currentCustomer)
		h.InsertToShopCart(c)

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.productId))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestGetListShopCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	LoadEnv()
	
	// Create database connection
	dbname := os.Getenv("DB_name")
	db, err := connectDB(dbname)
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}

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
	r := product.NewRepoProduct(db)
	s := product.NewService(r)
	h := handler.NewProductHandler(s)

	for _, testCase := range testCases {
		currentCustomer := customer.Customer{
			ID: testCase.customerId,
		}
		reqBody := url.Values{}
		reqBody.Set("id", testCase.id)
		g := gin.New()
		req := httptest.NewRequest(http.MethodGet, "listshopcart", strings.NewReader(reqBody.Encode()))
		res := httptest.NewRecorder()
		g.ServeHTTP(res, req)
		c, _ := gin.CreateTestContext(res)
		c.Set("currentCustomer", currentCustomer)
		h.GetListProductShopCart(c)

		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.productId))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}
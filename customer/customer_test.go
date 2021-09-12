package customer_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"graphql/customer"
	"graphql/handler"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

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

func TestMain(m *testing.M) {
	setup()
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

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

func setup() {
	// Create database connection
	LoadEnv()
	db, err := connectDB("")
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}

	dbname := os.Getenv("DB_name")

	// Cleaning database before testing
	db.Exec("DROP DATABASE IF EXISTS "+dbname)
	
	// Create test database
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	db.Close()
	
	// Create database connection
	db, err = connectDB(dbname)
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}

	// Create table
	db.MustExec(`CREATE TABLE customers(
							id int primary key,
							name varchar(30),
							email varchar(30),
							phone varchar(30),
							password varchar(30),
							salt varchar(30),
							avatar varchar(30),
							created_at datetime,
							updated_at datetime)`)
	// if err != nil {
	// 	log.Printf("Error %s when creating table", err)
	// 	return
	// }
	
	// Insert dummy data
	customer := customer.Customer{}
	customer.ID = 1
	customer.Name = "joni"
	customer.Email = "jon@j.j"
	customer.Password = "jojo"
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	_, err =db.Exec(`INSERT INTO customers
				(id, name, email, password, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?)`,
				customer.ID, customer.Name, customer.Email,
				customer.Password, customer.CreatedAt, customer.UpdatedAt)
	if err != nil {
		log.Printf("Error %s when inserting dummy data", err)
	}
}

func TestCreateCustomer(t *testing.T) {
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
		name string
		email string
		password string
		confirmPassword string
		expectCode int
		expectMsg string
	}{
		{
			testName: "success",
			name: "joni2",
			email: "jon@k.k",
			password: "jiji",
			confirmPassword: "jiji",
			expectCode: http.StatusOK,
			expectMsg: "account successfully created",
		},
		{
			testName: "fail used email",
			name: "joni3",
			email: "jon@k.k",
			password: "jojo",
			confirmPassword: "jojo",
			expectCode: http.StatusBadRequest,
			expectMsg: "email has been used",
		},
		{
			testName: "fail different password",
			name: "joni4",
			email: "jon@l.l",
			password: "jojo",
			confirmPassword: "jiji",
			expectCode: http.StatusBadRequest,
			expectMsg: "password and confirm password is different",
		},
	}

	// setting handler
	r := customer.NewRepo(db)
	s := customer.NewCustomerService(r)
	h := handler.NewHandlerCustomer(s, nil)

	for _, testCase := range testCases {
		reqBody, _ := json.Marshal(map[string]string{
			"name": testCase.name,
			"email": testCase.email,
			"password": testCase.password,
			"confirm_password": testCase.confirmPassword,
		})
		g := gin.New()
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		g.ServeHTTP(res, req)
		c, _ := gin.CreateTestContext(res)
		h.CreateCustomer(c)
		fmt.Println(res)
		
		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.name))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestLogin(t *testing.T) {
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
		email string
		password string
		expectCode int
		expectMsg string
	}{
		{
			testName: "success",
			email: "jon@k.k",
			password: "jiji",
			expectCode: http.StatusUnprocessableEntity,
			expectMsg: "success login",
		},
		{
			testName: "fail email not found",
			email: "jin@k.k",
			password: "jojo",
			expectCode: http.StatusUnprocessableEntity,
			expectMsg: "email not found",
		},
		{
			testName: "fail wrong password",
			email: "jon@k.k",
			password: "jojo",
			expectCode: http.StatusUnprocessableEntity,
			expectMsg: "different password",
		},
	}

	// setting handler
	r := customer.NewRepo(db)
	s := customer.NewCustomerService(r)
	h := handler.NewHandlerCustomer(s, nil)
	
	for _, testCase := range testCases {
		reqBody, _ := json.Marshal(map[string]string{
			"email": testCase.email,
			"password": testCase.password,
		})
		g := gin.New()
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		g.ServeHTTP(res, req)
		c, _ := gin.CreateTestContext(res)
		h.Login(c)
		
		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.email))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}

func TestUpdatePhoneCustomer(t *testing.T) {
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
		email string
		phone string
		expectCode int
		expectMsg string
	}{
		{
			testName: "success",
			email: "jon@k.k",
			phone: "0812345",
			expectCode: http.StatusUnprocessableEntity,
			expectMsg: "successfully udpate",
		},
		{
			testName: "fail email not found",
			email: "jin@k.k",
			phone: "0812345",
			expectCode: http.StatusUnprocessableEntity,
			expectMsg: "email not found",
		},
	}

	// setting handler
	r := customer.NewRepo(db)
	s := customer.NewCustomerService(r)
	h := handler.NewHandlerCustomer(s, nil)
	
	for _, testCase := range testCases {
		currentCustomer := customer.Customer{
			Email: testCase.email,
		}
		reqBody := url.Values{}
		reqBody.Set("phone", testCase.phone)
		g := gin.New()
		req := httptest.NewRequest(http.MethodPut, "/phone", strings.NewReader(reqBody.Encode()))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		g.ServeHTTP(res, req)
		c, _ := gin.CreateTestContext(res)
		c.Set("currentCustomer", currentCustomer)
		h.UpdatePhoneCustomer(c)
		
		assert.Equal(t, testCase.expectCode, res.Code)
		assert.True(t, strings.Contains(res.Body.String(), testCase.email))
		assert.True(t, strings.Contains(res.Body.String(), testCase.expectMsg))
	}
}


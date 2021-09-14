package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"olshop/auth"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

type customer struct {
	id       string
	email    string
	password string
}

// func Test_main(t *testing.T) {
// 	tests := []struct {
// 		name string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			main()
// 		})
// 	}
// }

func TestMiddleware(t *testing.T) {
	// TestCases list
	testCases := []struct {
		testName        string
		authType        string
		currentCustomer customer
	}{
		{
			testName: "success",
			authType: "Bearer",
			currentCustomer: customer{
				id:       "1",
				email:    "jojo",
				password: "jon",
			},
		},
		{
			testName: "fail",
			authType: "Bearer",
			currentCustomer: customer{
				id: "2",
				email: "jiji",
				password: "jin",
			},
		},
	}

	for _, testCase := range testCases {
		id, _ := strconv.Atoi(testCase.currentCustomer.id)
		token, err := auth.NewService().GenerateToken(id)
		assert.Nil(t, err)

		reqBody, _ := json.Marshal(map[string]string{
			"id" : testCase.currentCustomer.id,
			"email" : testCase.currentCustomer.email,
			"password" : testCase.currentCustomer.password,
		})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		c, r := gin.CreateTestContext(res)
		c.Request = req
		c.Request.Header.Set("Authorization", testCase.authType+" "+token)

		r.ServeHTTP(res, req)

		header := c.GetHeader("Authorization")
		assert.True(t, strings.Contains(header, testCase.authType))

		tokenSplit := strings.Split(header, " ")[1]

		tokenString, err := auth.NewService().ValidateToken(tokenSplit)
		assert.Nil(t, err)

		claim, ok := tokenString.Claims.(jwt.MapClaims)
		assert.True(t, ok)

		customerId := int(claim["customer_id"].(float64))
		assert.Equal(t, id, customerId)
	}
}

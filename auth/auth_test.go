package auth_test

import (
	"olshop/auth"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateValidateToken(t *testing.T) {
	// Testcases list
	testCases := []struct {
		testName string
		id       uint
		exp      int64
	}{
		{
			testName: "success",
			id:       1,
			exp:      time.Now().Add(time.Hour * 8).Unix(),
		},
		{
			testName: "success",
			id:       2,
			exp:      time.Now().Add(time.Hour * 8).Unix(),
		}, {
			testName: "success",
			id:       30,
			exp:      time.Now().Add(time.Hour * 8).Unix(),
		},
	}

	for _, testCase := range testCases {
		newToken, err := auth.NewService("coba", "cobalagi").GenerateToken(testCase.id)
		assert.Nil(t, err)

		newToken1, err := auth.NewService("coba", "cobalagi").GenerateTokenSeller(testCase.id)
		assert.Nil(t, err)

		token, er := auth.NewService("coba", "cobalagi").ValidateToken(newToken)
		assert.Nil(t, er)

		token1, err := auth.NewService("coba", "cobalagi").ValidateTokenSeller(newToken1)
		assert.NoError(t, err)

		claim, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok)

		claim1, ok := token1.Claims.(jwt.MapClaims)
		assert.True(t, ok)

		id := int(claim["customer_id"].(float64))
		time := int64(claim["exp"].(float64))
		assert.Equal(t, testCase.id, id)
		assert.Equal(t, testCase.exp, time)

		id1 := int(claim1["seller_id"].(float64))
		time1 := int64(claim1["exp"].(float64))
		assert.Equal(t, testCase.id, id1)
		assert.Equal(t, testCase.exp, time1)
	}
}

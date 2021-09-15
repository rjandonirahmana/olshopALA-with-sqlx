package auth_test

import (
	"olshop/auth"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateValidateToken(t *testing.T) {
	// Testcases list
	testCases := []struct{
		testName string
		id int
	}{
		{
			testName: "success",
			id: 1,
		},
		{
			testName: "success",
			id: 2,
		},
	}

	for _, testCase := range testCases {
		newToken, err := auth.NewService().GenerateToken(testCase.id)
		assert.Nil(t, err)

		token, er := auth.NewService().ValidateToken(newToken)
		assert.Nil(t, er)

		claim, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok)

		id := int(claim["customer_id"].(float64))
		assert.Equal(t, testCase.id, id)
	}
}
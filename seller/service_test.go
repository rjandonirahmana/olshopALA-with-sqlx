package seller

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSellerService(t *testing.T) {
	testCases := []struct {
		name   string
		seller InputSeller
		err    error
	}{
		{
			name:   "test1",
			seller: InputSeller{Name: "coba", Email: "coba lagi", Password: "ngasal"},
			err:    nil,
		}, {
			name:   "test1",
			seller: InputSeller{Name: "coba", Email: "aku1", Password: "12345"},
			err:    nil,
		},
	}

	for _, test := range testCases {
		seller, err := services.Register(test.seller)
		assert.Equal(t, test.err, err)
		fmt.Println(seller)
	}
}

func TestLogin(t *testing.T) {
	testCases := []struct {
		name  string
		input InputLoginSeller
		err   error
	}{
		{
			name:  "test1",
			input: InputLoginSeller{Email: "aku1", Password: "12345"},
			err:   nil,
		}, {
			name:  "test1",
			input: InputLoginSeller{Email: "coba lagi", Password: "ngasal"},
			err:   nil,
		},
	}

	for _, test := range testCases {
		seller, err := services.LoginSeller(test.input)
		assert.Equal(t, test.err, err)

		fmt.Println(seller)
	}
}

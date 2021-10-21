package seller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"olshop/auth"
	"olshop/seller"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupdb() *sqlx.DB {

	dbUserName := "---------"
	dbName := "--------"
	dbPass := "----------"

	dbString := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", dbUserName, dbPass, dbName)
	db, err := sqlx.Connect("postgres", dbString)

	if err != nil {
		panic(err)
	}

	return db
}

var (
	db            = setupdb()
	repo          = seller.NewRepository(db)
	service       = seller.NewService(repo, "ngasal")
	authorization = auth.NewService("coba", "cobalagi")
	handlerseller = NewHandlerSeller(service, authorization)
)

func TestRegister(t *testing.T) {
	testCases := []struct {
		name   string
		seller seller.InputSeller
		err    error
	}{
		{
			name:   "test1",
			seller: seller.InputSeller{Name: "cobalah", Email: "rjan@gmail.com", Password: "12345678910", ConfirmPassword: "12345678910"},
			err:    nil,
		}, {
			name:   "test1",
			seller: seller.InputSeller{Name: "coba", Email: "aku1", Password: "12345", ConfirmPassword: "12345"},
			err:    nil,
		},
	}

	for _, test := range testCases {
		reqBody, err := json.Marshal(test.seller)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = req

		handlerseller.Register(c)
		fmt.Println(res)
	}
}

package customer

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupdb() *sqlx.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading environment variables")
		os.Exit(1)
	}

	dbUserName := os.Getenv("DB_username")
	dbName := os.Getenv("DB_name")
	dbPass := os.Getenv("DB_password")

	dbString := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", dbUserName, dbPass, dbName)
	db, err := sqlx.Connect("postgres", dbString)

	if err != nil {
		panic(err)
	}

	return db
}

var (
	db      = setupdb()
	repo    = NewRepo(db)
	service = NewCustomerService(repo, "abc")
)

func TestRegister(t *testing.T) {
	testCases := []struct {
		codetest uint
		request  Customer
		err      error
	}{
		{
			codetest: 1,
			request:  Customer{Name: "ngasal", Email: "rjandoni", Password: "apa", Phone: "12345", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			err:      fmt.Errorf("email has been used by another"),
		}, {
			codetest: 1,
			request:  Customer{Name: "ngasal", Email: "rahp", Password: "apa", Phone: "12345", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			err:      fmt.Errorf("email has been used by another"),
		},
	}

	for _, test := range testCases {
		customer, err := service.Register(test.request)

		assert.Equal(t, test.err, err)

		fmt.Println(customer)

	}
}

func TestLogin(t *testing.T) {
	testCases := []struct {
		name  string
		input InputLogin
		err   error
	}{
		{
			name:  "test1",
			input: InputLogin{Email: "rjandoni", Password: "12345"},
			err:   nil,
		}, {
			name:  "test1",
			input: InputLogin{Email: "rahp", Password: "apa"},
			err:   nil,
		},
	}

	for _, test := range testCases {
		seller, err := service.LoginCustomer(test.input)
		assert.Equal(t, test.err, err)

		fmt.Println(seller)
	}
}

func TestChangePassword(t *testing.T) {
	testCases := []struct {
		name        string
		id          uint
		oldpassword string
		newpassword string
		err         error
	}{
		{
			name:        "test1",
			id:          12,
			oldpassword: "12345",
			newpassword: "apa",
			err:         nil,
		},
	}

	for _, test := range testCases {
		seller, err := service.ChangePassword(test.oldpassword, test.newpassword, test.id)
		assert.Equal(t, test.err, err)

		fmt.Println(seller)
	}

}
func TestIsemailAvailable(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		err   error
	}{
		{
			name:  "1",
			email: "ngasal",
			err:   fmt.Errorf("email has been used by another"),
		},
		{
			name:  "1",
			email: "doniiiaaa@gmail.com",
			err:   nil,
		}, {
			name:  "2",
			email: "ngasal",
			err:   fmt.Errorf("email has been used by another"),
		}, {
			name:  "2",
			email: "coba ada error ga",
			err:   nil,
		}, {
			name:  "2",
			email: "coba1",
			err:   fmt.Errorf("email has been used by another"),
		},
	}

	for _, test := range testCases {
		err := repo.IsEmailAvailable(test.email)
		fmt.Printf("testcaseeeeeeee   ============================ %s=======================", test.name)
		fmt.Println(err)
		assert.Equal(t, test.err, err)
	}
}

func TestGetbyEmail(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		err   error
	}{
		{
			name:  "test1",
			email: "ngasal",
			err:   nil,
		}, {
			name:  "2",
			email: "ATENG",
			err:   nil,
		},
	}

	for _, test := range testCases {
		user, err := repo.GetCustomerByEmail(test.email)
		assert.NoError(t, err)

		fmt.Println(user)
	}
}

func TestUpdateAvatar(t *testing.T) {
	testCases := []struct {
		file string
		id   uint
		err  error
	}{
		{
			file: "dgddtddtdtdd",
			id:   3,
			err:  nil,
		},
	}

	for _, test := range testCases {
		err := repo.ChangeAvatar(test.file, test.id)
		assert.Equal(t, test.err, err)
	}
}

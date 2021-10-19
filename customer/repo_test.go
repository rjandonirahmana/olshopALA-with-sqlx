package customer

import (
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupdb() *sqlx.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading environment variables")
		os.Exit(1)
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
		os.Getenv("DB_username"),
		os.Getenv("DB_password"),
		os.Getenv("DB_host"),
		os.Getenv("DB_port"),
		os.Getenv("DB_name"))
	db, err := sqlx.Connect("mysql", conn)
	if err != nil {
		panic(err)
	}

	return db
}

var (
	db   = setupdb()
	repo = NewRepo(db)
	// service = NewCustomerService(repo)
)

func TestIsemailAvailable(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		err   error
	}{
		{
			name:  "1",
			email: "doniiiaaa@gmail.com",
			err:   nil,
		}, {
			name:  "2",
			email: "siapaaja",
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
			email: "baruyaaa@gmail.com",
			err:   nil,
		}, {
			name:  "2",
			email: "123456",
			err:   nil,
		},
	}

	for _, test := range testCases {
		user, err := repo.GetCustomerByEmail(test.email)
		assert.NoError(t, err)

		fmt.Println(user)
	}
}

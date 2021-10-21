package seller

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

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
	db       = setupdb()
	repo     = NewRepository(db)
	services = NewService(repo, "ngasal")

	// service = NewCustomerService(repo)
)

func TestGetbyEmail(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		err   error
	}{
		{
			name:  "test1",
			email: "coba lagi",
			err:   nil,
		}, {
			name:  "2",
			email: "nyoba",
			err:   fmt.Errorf("email not found"),
		}, {
			name:  "2",
			email: "aku1",
			err:   nil,
		},
	}

	for _, test := range testCases {
		user, err := repo.GetSellerByEmail(test.email)
		assert.Equal(t, test.err, err)
		fmt.Println(user)

	}
}

func TestIsEmailAvalaible(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		err   error
	}{
		{
			name:  "test1",
			email: "aku1",
			err:   fmt.Errorf("email has been used by another user"),
		}, {
			name:  "2",
			email: "123456",
			err:   nil,
		}, {
			name:  "3",
			email: "coba lagi",
			err:   fmt.Errorf("email has been used by another user"),
		}, {
			name:  "4",
			email: ".,cgd.mc",
			err:   nil,
		},
	}

	for _, test := range testCases {
		err := repo.IsEmailAvailable(test.email)

		assert.Equal(t, test.err, err)
		fmt.Println(err)
	}
}

// func TestIsCreateSeller(t *testing.T) {
// 	testCases := []struct {
// 		name   string
// 		seller Seller
// 		err    error
// 	}{
// 		{
// 			name:   "test1",
// 			seller: Seller{Name: "coba", Email: "nyoba", Password: "ngasal", Salt: "ngasal", CreatedAt: time.Now(), UpdatedAt: time.Now()},
// 			err:    nil,
// 		},
// 	}

// 	for _, test := range testCases {
// 		seller, err := repo.CreateSeller(test.seller)
// 		assert.NoError(t, err)
// 		fmt.Println(seller)
// 	}
// }

func TestGetByID(t *testing.T) {
	testCases := []struct {
		name string
		id   uint
		err  error
	}{
		{
			name: "test1",
			id:   31,
			err:  nil,
		}, {
			name: "test2",
			id:   32,
			err:  nil,
		}, {
			name: "test3",
			id:   1000,
			err:  sql.ErrNoRows,
		},
	}

	for _, test := range testCases {
		seller, err := repo.GetSellerByID(test.id)
		assert.Equal(t, test.err, err)

		fmt.Println(seller)
	}
}

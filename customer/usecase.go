package customer

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ServiceCustomer struct {
	repo Repository
}

type CustomerInt interface {
	Register(customer Customer) (Customer, error)
	LoginCustomer(input InputLogin) (Customer, error)
}

func NewCustomerService(repository Repository) *ServiceCustomer {
	return &ServiceCustomer{repo: repository}
}

func (s *ServiceCustomer) Register(customer Customer) (Customer, error) {

	customer.Salt = "salt"
	customer.Password += customer.Salt
	customer.CreatedAt = time.Now()

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.MinCost)
	if err != nil {
		return Customer{}, err
	}

	customer.Password = string(hashpassword)
	id, _ := s.repo.GetLastID()
	customer.ID = id + 1

	customer, err = s.repo.RegisterUser(customer)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (s *ServiceCustomer) LoginCustomer(input InputLogin) (Customer, error) {

	salt := "salt"
	password := input.Password + salt

	customer, err := s.repo.GetCustomerByEmail(input.Email)
	if err != nil {
		return Customer{}, errors.New("email not found")

	}

	err1 := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password))
	if err1 != nil {
		return Customer{}, errors.New("different password")
	}

	return customer, nil

}

func (s *ServiceCustomer) UpdateCustomerPhone(phone int64, email string) error {

	_, err := s.repo.GetCustomerByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}

	err = s.repo.UpdateCustomerPhone(email, phone)
	if err != nil {
		return err
	}

	return nil
}

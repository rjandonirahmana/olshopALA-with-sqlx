package customer

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type ServiceCustomer struct {
	repo Repository
}

type CustomerInt interface {
	Register(customer Customer) error
	LoginCustomer(input InputLogin) (Customer, error)
}

func NewCustomerService(repository Repository) *ServiceCustomer {
	return &ServiceCustomer{repo: repository}
}

func (s *ServiceCustomer) Register(customer Customer) error {

	customer.Salt = "salt"
	customer.Password += customer.Salt

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	customer.Password = string(hashpassword)

	err = s.repo.RegisterUser(customer)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *ServiceCustomer) LoginCustomer(input InputLogin) (Customer, error) {

	salt := "salt"
	password := input.Password + salt

	customer, err := s.repo.GetCustomerByEmail(input.Email)
	if err {
		return Customer{}, errors.New("email not found")

	}

	err1 := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password))
	if err1 != nil {
		return Customer{}, errors.New("this isnt your account")
	}

	return customer, nil

}

package customer

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type ServiceCustomer struct {
	repo Repository
}

type CustomerInt interface {
	Register(customer Customer) error
}

func NewCustomerService(repository Repository) *ServiceCustomer {
	return &ServiceCustomer{repo: repository}
}

func (s *ServiceCustomer) Register(customer Customer) error {

	customer.Salt = "salt"
	customer.Password += customer.Salt

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.MaxCost)
	if err != nil {
		return err
	}

	customer.Password = string(hashpassword)

	err = s.repo.RegisterUser(customer)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceCustomer) LoginCustomer(input Customer) error {

	input.Salt = "salt"
	password := input.Password + input.Salt

	customer, err := s.repo.GetCustomerByEmail(input.Email)
	if err {
		return errors.New("email not found")

	}

	err1 := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password))
	if err1 != nil {
		return errors.New("this isnt your account")
	}

	return nil

}

func (s *ServiceCustomer) IsEmailAvailable(input Customer) bool {

	_, available := s.repo.GetCustomerByEmail(input.Email)
	if !available {
		return true
	}

	return false

}

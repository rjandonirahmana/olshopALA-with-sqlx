package usecase

import (
	"graphql/repo"

	"golang.org/x/crypto/bcrypt"
)

type ServiceCustomer struct {
	repo repo.Repository
}

type CustomerInt interface {
	Register(customer repo.Customer) error
}

func NewCustomerService(repository repo.Repository) *ServiceCustomer {
	return &ServiceCustomer{repo: repository}
}

func (s *ServiceCustomer) Register(customer repo.Customer) error {

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.MaxCost)
	if err != nil {
		return err
	}

	customer.Password = string(hashpassword) + customer.Salt

	err = s.repo.RegisterUser(customer)
	if err != nil {
		return err
	}

	return nil
}

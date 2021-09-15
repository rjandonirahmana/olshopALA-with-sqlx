package customer

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/crypto/bcrypt"
)

const AllowedExtensions = ".jpeg,.jpg"
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ServiceCustomer struct {
	repo Repository
}

type CustomerInt interface {
	Register(customer Customer) (Customer, error)
	LoginCustomer(input InputLogin) (Customer, error)
	UpdateCustomerPhone(phone string, email string) error
	GetCustomerByID(id int) (Customer, error)
	IsEmailAvailable(email string) (bool, error)
	ChangeProfile(profile []byte, name string, id int) (Customer, error)
	ChangePassword(oldpassword, newPassword string, id int) (Customer, error)
	DeleteCustomer(id int, password string) error
}

func NewCustomerService(repo Repository) *ServiceCustomer {
	return &ServiceCustomer{repo: repo}
}

func (s *ServiceCustomer) Register(customer Customer) (Customer, error) {

	_, err := s.IsEmailAvailable(customer.Email)
	if err != nil {
		return Customer{}, errors.New("email has been used")
	}

	customer.Salt = RandStringBytes(len(customer.Password) + 9)
	customer.Password += customer.Salt
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.MinCost)
	if err != nil {
		return Customer{}, err
	}

	customer.Password = string(hashpassword)

	customer, err = s.repo.RegisterUser(customer)
	if err != nil {
		return Customer{}, err
	}

	customer, err = s.repo.GetCustomerByEmail(customer.Email)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (s *ServiceCustomer) LoginCustomer(input InputLogin) (Customer, error) {

	customer, err := s.repo.GetCustomerByEmail(input.Email)
	if err != nil {
		return Customer{}, errors.New("email not found")

	}

	password := input.Password + customer.Salt

	err1 := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password))
	if err1 != nil {
		return Customer{}, errors.New("different password")
	}

	return customer, nil

}

func (s *ServiceCustomer) UpdateCustomerPhone(phone string, email string) error {

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

func (s *ServiceCustomer) GetCustomerByID(id int) (Customer, error) {

	customer, err := s.repo.GetCustomerByID(id)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (s *ServiceCustomer) IsEmailAvailable(email string) (bool, error) {
	customer, _ := s.repo.GetCustomerByEmail(email)

	if customer.Email == email {
		return false, errors.New("email has been used")
	}
	return true, nil
}

func (s *ServiceCustomer) ChangeProfile(profile []byte, name string, id int) (Customer, error) {

	mime := mimetype.Detect(profile)
	if strings.Index(AllowedExtensions, mime.Extension()) == -1 {
		return Customer{}, errors.New("File Type is not allowed, file type: " + mime.Extension())
	}

	profilesave := fmt.Sprintf("image/profile/%s,%s", name, mime.Extension())
	err := s.repo.ChangeAvatar(profilesave, id)
	if err != nil {
		return Customer{}, err
	}
	customer, err := s.repo.GetCustomerByID(id)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *ServiceCustomer) ChangePassword(oldpassword, newPassword string, id int) (Customer, error) {
	customer, err := s.repo.GetCustomerByID(id)
	if err != nil {
		return Customer{}, err
	}
	oldpassword += customer.Salt

	err1 := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(oldpassword))
	if err1 != nil {
		return Customer{}, errors.New("please input your password correctly")
	}
	newPassword += customer.Salt

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
	if err != nil {
		return Customer{}, err
	}

	err = s.repo.ChangePassword(string(hashpassword), id)
	if err != nil {
		return Customer{}, err
	}

	customer, err = s.repo.GetCustomerByID(id)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil

}

func (s *ServiceCustomer) DeleteCustomer(id int, password string) error {
	customer, err := s.repo.GetCustomerByID(id)
	if err != nil {
		return err
	}

	password += customer.Salt
	err1 := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password))
	if err1 != nil {
		return errors.New("cant delete your account if you dont know your password")
	}

	err = s.repo.DeleteCustomer(id)
	if err != nil {
		return err
	}

	return nil
}

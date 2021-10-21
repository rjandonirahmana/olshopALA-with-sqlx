package customer

import (
	"crypto/sha256"
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
	repo   Repository
	secret string
}

type CustomerInt interface {
	Register(customer Customer) (Customer, error)
	LoginCustomer(input InputLogin) (Customer, error)
	UpdateCustomerPhone(phone string, id uint) error
	GetCustomerByID(id uint) (Customer, error)
	ChangeProfile(profile []byte, name string, id uint) (Customer, error)
	ChangePassword(oldpassword, newPassword string, id uint) (Customer, error)
	DeleteCustomer(id uint, password string) error
}

func NewCustomerService(repo Repository, secret string) *ServiceCustomer {
	return &ServiceCustomer{repo: repo, secret: secret}
}

func (s *ServiceCustomer) Register(customer Customer) (Customer, error) {

	err := s.repo.IsEmailAvailable(customer.Email)
	if err != nil {
		return Customer{}, err
	}

	customer.Salt = RandStringBytes(len(customer.Password) + 9)
	customer.Password += customer.Salt
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	h := sha256.New()
	h.Write([]byte(customer.Password))
	customer.Password = fmt.Sprintf("%x", h.Sum([]byte(s.secret)))

	customer, err = s.repo.RegisterUser(customer)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (s *ServiceCustomer) LoginCustomer(input InputLogin) (Customer, error) {

	customer, err := s.repo.GetCustomerByEmail(input.Email)
	if err != nil {
		return Customer{}, err

	}
	input.Password += customer.Salt
	h := sha256.New()
	h.Write([]byte(input.Password))
	hashpassword := fmt.Sprintf("%x", h.Sum([]byte(s.secret)))

	if customer.Password != hashpassword {
		return Customer{}, errors.New("your input password false, please input your password correctly")
	}

	return customer, nil

}

func (s *ServiceCustomer) UpdateCustomerPhone(phone string, id uint) error {

	err := s.repo.UpdateCustomerPhone(id, phone)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceCustomer) GetCustomerByID(id uint) (Customer, error) {

	customer, err := s.repo.GetCustomerByID(id)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (s *ServiceCustomer) ChangeProfile(profile []byte, name string, id uint) (Customer, error) {

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

func (s *ServiceCustomer) ChangePassword(oldpassword, newPassword string, id uint) (Customer, error) {
	customer, err := s.repo.GetCustomerByID(id)
	if err != nil {
		return Customer{}, err
	}
	oldpassword += customer.Salt
	h := sha256.New()
	h.Write([]byte(oldpassword))
	hashpassword := fmt.Sprintf("%x", h.Sum([]byte(s.secret)))

	if customer.Password != hashpassword {
		return Customer{}, errors.New("your input password false, cannot change password")
	}

	newPassword += customer.Salt
	h = sha256.New()
	h.Write([]byte(newPassword))
	hashpassword = fmt.Sprintf("%x", h.Sum([]byte(s.secret)))

	err = s.repo.ChangePassword(hashpassword, id)
	if err != nil {
		return Customer{}, err
	}

	customer.Password = hashpassword

	return customer, nil

}

func (s *ServiceCustomer) DeleteCustomer(id uint, password string) error {
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

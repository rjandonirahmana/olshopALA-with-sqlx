package seller

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type service struct {
	repo   RepoInt
	secret string
}

type Service interface {
	Register(input InputSeller) (Seller, error)
	GetSellerByID(id uint) (Seller, error)
	LoginSeller(input InputLoginSeller) (Seller, error)
}

func NewService(repo RepoInt, secret string) *service {
	return &service{repo: repo, secret: secret}
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *service) Register(input InputSeller) (Seller, error) {
	err := s.repo.IsEmailAvailable(input.Email)
	if err != nil {
		return Seller{}, err
	}

	salt := RandStringBytes(len(s.secret))
	input.Password += salt

	h := sha256.New()
	h.Write([]byte(input.Password))

	Password := fmt.Sprintf("%x", h.Sum([]byte(s.secret)))

	seller := Seller{
		Name:      input.Name,
		Email:     input.Email,
		Password:  Password,
		Salt:      salt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	seller, err = s.repo.CreateSeller(seller)
	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

func (s *service) LoginSeller(input InputLoginSeller) (Seller, error) {

	seller, err := s.repo.GetSellerByEmail(input.Email)
	if err != nil {
		return Seller{}, err

	}
	input.Password += seller.Salt
	h := sha256.New()
	h.Write([]byte(input.Password))
	hashpassword := fmt.Sprintf("%x", h.Sum([]byte(s.secret)))

	if seller.Password != hashpassword {
		fmt.Println(hashpassword)
		return Seller{}, errors.New("your input password false, please input your password correctly")
	}

	return seller, nil

}

func (s *service) GetSellerByID(id uint) (Seller, error) {
	seller, err := s.repo.GetSellerByID(id)
	if err != nil {
		return seller, err
	}

	return seller, nil
}

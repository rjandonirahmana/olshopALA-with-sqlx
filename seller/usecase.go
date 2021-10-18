package seller

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo RepoInt
}

type Service interface {
	Register(seller Seller) (Seller, error)
	GetSellerByID(id int) (Seller, error)
}

func NewService(repo RepoInt) *service {
	return &service{repo: repo}
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *service) Register(seller Seller) (Seller, error) {

	seller.Salt = RandStringBytes(len(seller.Password))
	seller.Password += seller.Salt
	seller.CreatedAt = time.Now()
	seller.UpdatedAt = time.Now()

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(seller.Password), bcrypt.MinCost)
	if err != nil {
		return Seller{}, err
	}

	seller.Password = string(hashpassword)

	seller, err = s.repo.CreateSeller(seller)
	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

func (s *service) GetSellerByID(id int) (Seller, error) {
	seller, err := s.repo.GetSellerByID(id)
	if err != nil {
		return seller, err
	}

	return seller, nil
}

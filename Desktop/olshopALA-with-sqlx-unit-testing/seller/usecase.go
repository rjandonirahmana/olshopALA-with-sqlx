package seller

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	service RepoInt
}

func (s *Service) Register(seller Seller) (Seller, error) {

	seller.Salt = "salt"
	seller.Password += seller.Salt
	seller.CreatedAt = time.Now()
	seller.UpdatedAt = time.Now()

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(seller.Password), bcrypt.MinCost)
	if err != nil {
		return Seller{}, err
	}

	seller.Password = string(hashpassword)

	id, _ := s.service.GetLastID()
	seller.ID = id + 1

	seller, err = s.service.CreateSeller(seller)
	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

func (s *Service) LoginSeller()

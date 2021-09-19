package seller

import (
	"fmt"
	"olshop/customer"
	"olshop/product"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	service RepoInt
	product product.RepoProduct
}

type ServiceInt interface {
	Register(seller Seller) (Seller, error)
	GetSellerById(id int) (Seller, error)
	Login(input InputLogin) (Seller, error)
	InsertProduct(product product.Product, type_name string) (product.Product, error)
}

func NewServiceSeller(repo RepoInt, repoProduct product.RepoProduct) *Service {
	return &Service{service: repo, product: repoProduct}
}

func (s *Service) Register(seller Seller) (Seller, error) {

	available, err := s.service.IsEmailAvailable(seller.Email)
	if !available || err != nil {

		return Seller{}, fmt.Errorf("email has been used or error : %v", err)
	}

	seller.Salt = customer.RandStringBytes(len(seller.Password) + 10)
	seller.Password += seller.Salt
	seller.CreatedAt = time.Now()
	seller.UpdatedAt = time.Now()

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(seller.Password), bcrypt.MinCost)
	if err != nil {
		return Seller{}, err
	}

	seller.Password = string(hashpassword)

	seller, err = s.service.CreateSeller(seller)
	if err != nil {
		fmt.Println("disini masuk ga?")
		return Seller{}, err
	}

	updatedSeller, err := s.service.GetSellerByEmail(seller.Email)
	if err != nil {
		return Seller{}, err
	}

	return updatedSeller, nil
}

func (s *Service) Login(input InputLogin) (Seller, error) {
	seller, err := s.service.GetSellerByEmail(input.Email)
	if err != nil {
		return Seller{}, fmt.Errorf("email not found or error : %v", err)
	}

	input.Password += seller.Salt

	err = bcrypt.CompareHashAndPassword([]byte(seller.Password), []byte(input.Password))
	if err != nil {
		return Seller{}, fmt.Errorf("cant login, please input your password correctly")
	}

	return seller, nil
}

func (s *Service) GetSellerById(id int) (Seller, error) {
	seller, err := s.service.GetSellerByID(id)
	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

func (s *Service) InsertProduct(product product.Product, type_name string) (product.Product, error) {
	id, err := s.service.SelectTypeProductID(type_name)
	if err != nil {
		return product, err
	}
	product.Category_id = id

	ID, err := s.service.InsertProduct(product)
	if err != nil {
		return product, err
	}

	products, err := s.product.GetProductByID(int(ID))
	if err != nil {
		return product, err
	}

	return products, nil

}

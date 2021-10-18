package product

import "errors"

type serviceProduct struct {
	repo RepoProduct
}

type ServiceProductInt interface {
	GetProductCategory(id int) ([]Product, error)
	GetProductByid(id int) (Products, error)
}

func NewService(repo RepoProduct) *serviceProduct {
	return &serviceProduct{repo: repo}
}

func (s *serviceProduct) GetProductCategory(id int) ([]Product, error) {

	products, err := s.repo.GetByCategoryID(id)
	if err != nil {
		return []Product{}, err
	}

	if len(products) == 0 {
		return []Product{}, errors.New("products is not found")
	}

	return products, nil
}

func (s *serviceProduct) GetProductByid(id int) (Products, error) {
	product, err := s.repo.GetProductByID(id)

	if err != nil {
		return Products{}, err
	}

	return product, nil
}

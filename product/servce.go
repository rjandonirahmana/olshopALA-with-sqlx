package product

import "errors"

type serviceProduct struct {
	repo RepoProduct
}

type ServiceProductInt interface {
	GetProductCategory(name_category string) ([]Product, error)
	AddShoppingCart(customerid int) (int, error)
	InsertProductByCartID(customerid, productid, cartid int) (Product, error)
	GetListShopCart(cartid int, customerid int) ([]Product, error)
	DeleteListOnshoppingCart(cartid, customerid, productid int) ([]Product, error)
}

func NewService(repo RepoProduct) *serviceProduct {
	return &serviceProduct{repo: repo}
}

func (s *serviceProduct) GetProductCategory(name_category string) ([]Product, error) {

	products, err := s.repo.GetProductByCategoryName(name_category)
	if err != nil {
		return []Product{}, err
	}

	if len(products) == 0 {
		return []Product{}, errors.New("products is not found")
	}

	return products, nil
}

func (s *serviceProduct) AddShoppingCart(customerid int) (int, error) {
	id, err := s.repo.GetLastID()
	if err != nil {
		return 0, err
	}

	id += 1

	err = s.repo.CreateCart(customerid, id)
	if err != nil {
		return 0, err
	}

	return id, nil

}

func (s *serviceProduct) InsertProductByCartID(customerid, productid, cartid int) (Product, error) {

	product, err := s.repo.GetProductByID(productid)
	if err != nil {
		return Product{}, err
	}

	cart, err := s.repo.GetShopCartIDCustomer(customerid, cartid)
	if err != nil {
		return Product{}, err
	}

	err = s.repo.InsertShoppingCart(cart.ID, productid, product.Price, product.Name)
	if err != nil {
		return Product{}, err

	}

	return product, nil

}

func (s *serviceProduct) GetListShopCart(cartid int, customerid int) ([]Product, error) {

	products, err := s.repo.GetListCartByID(cartid, customerid)
	if err != nil {
		return []Product{}, err
	}

	return products, nil
}

func (s *serviceProduct) DeleteListOnshoppingCart(cartid, customerid, productid int) ([]Product, error) {
	err := s.repo.DeleteProductInShopCart(cartid, customerid, productid)
	if err != nil {
		return []Product{}, err
	}

	productLeft, err := s.GetListShopCart(cartid, customerid)
	if err != nil {
		return []Product{}, err
	}

	return productLeft, nil
}

package product

type serviceProduct struct {
	repo RepoProduct
}

type ServiceProductInt interface {
	GetProductCategory(name_category string) ([]Product, error)
	AddShoppingCart(customerid int) error
	InsertProductByCartID(productid, cartid int) (Product, error)
	GetListShopCart(cartid int, customerid int) ([]Product, error)
}

func NewService(repo RepoProduct) *serviceProduct {
	return &serviceProduct{repo: repo}
}

func (s *serviceProduct) GetProductCategory(name_category string) ([]Product, error) {

	products, err := s.repo.GetProductByCategoryName(name_category)
	if err != nil {
		return []Product{}, err
	}

	return products, nil
}

func (s *serviceProduct) AddShoppingCart(customerid int) error {
	id, err := s.repo.GetLastID()
	if err != nil {
		return err
	}

	id += 1

	err = s.repo.CreateCart(customerid, id)
	if err != nil {
		return err
	}

	return nil

}

func (s *serviceProduct) InsertProductByCartID(productid, cartid int) (Product, error) {

	product, err := s.repo.GetProductByID(productid)
	if err != nil {
		return Product{}, err
	}

	err = s.repo.InsertShoppingCart(cartid, productid, product.Price, product.Name)
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

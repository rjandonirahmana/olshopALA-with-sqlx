package shopcart

type service struct {
	repository Repository
}

type Service interface {
	GetListInShopCart(cartid int, customerid int) ([]ShopeCart, error)
	DeleteListOnshoppingCart(cartid, customerid, productid int) ([]ShopeCart, error)
	GetShopCartCustomer(customerid int) ([]Cart, error)
	DecreaseProductShopCart(customerid, productid, cartid int) ([]ShopeCart, error)
}

func NewService(repo Repository) *service {
	return &service{repository: repo}
}

func (s *service) GetListInShopCart(cartid int, customerid int) ([]ShopeCart, error) {

	_, err := s.repository.GetShopCartIDCustomer(customerid, cartid)
	if err != nil {
		return []ShopeCart{}, err
	}

	shopcart, err := s.repository.GetListCartByID(cartid)
	if err == nil && len(shopcart) == 0 {
		return []ShopeCart{}, nil
	}
	if err != nil {
		return []ShopeCart{}, err
	}

	return shopcart, nil
}

func (s *service) DeleteListOnshoppingCart(cartid, customerid, productid int) ([]ShopeCart, error) {
	err := s.repository.DeleteProductInShopCart(cartid, customerid, productid)
	if err != nil {
		return []ShopeCart{}, err
	}

	productLeft, err := s.repository.GetListCartByID(cartid)
	if err != nil {
		return []ShopeCart{}, err
	}

	return productLeft, nil
}

func (s *service) GetShopCartCustomer(customerid int) ([]Cart, error) {

	cart, err := s.repository.ShopCartCustomer(customerid)
	if err != nil {
		return []Cart{}, err
	}

	return cart, nil
}

func (s *service) DecreaseProductShopCart(customerid, productid, cartid int) ([]ShopeCart, error) {
	_, err := s.repository.GetShopCartIDCustomer(customerid, cartid)
	if err != nil {
		return []ShopeCart{}, err
	}

	err = s.repository.DecreaseQuantitInShopCart(cartid, productid)
	if err != nil {
		return []ShopeCart{}, err
	}

	err = s.repository.DeleteAllWhenQuantity0()
	if err != nil {
		return []ShopeCart{}, err
	}

	productLeft, err := s.repository.GetListCartByID(cartid)
	if err != nil {
		return []ShopeCart{}, err
	}

	return productLeft, nil

}

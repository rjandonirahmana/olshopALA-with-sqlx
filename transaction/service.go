package transaction

import (
	"graphql/product"
	"time"
)

type ServiceTrans struct {
	repo        RepoTransaction
	repoProduct product.RepoProduct
}

type ServiceTransInt interface {
	CreateTransaction(idProduct int, customerID, quantity int) (Transactions, error)
}

func (s *ServiceTrans) CreateTransaction(idProduct, customerID, quantity int) (Transactions, error) {
	idTrans, _ := s.repo.GetLastTransaction()

	product, err := s.repoProduct.GetProductByID(idProduct)
	if err != nil {
		return Transactions{}, err
	}
	t := Transactions{}
	t.ID = idTrans + 1
	t.ID_product = product.ID
	t.Price = product.Price
	t.CreatedAt = time.Now()
	t.CustomerID = customerID
	t.Quantity = quantity

	err = s.repo.InserTransaction(t)
	if err != nil {
		return Transactions{}, err
	}

	trans, err := s.repo.GetDetailTransaction(t.ID)
	if err != nil {
		return Transactions{}, err
	}

	return trans, nil

}

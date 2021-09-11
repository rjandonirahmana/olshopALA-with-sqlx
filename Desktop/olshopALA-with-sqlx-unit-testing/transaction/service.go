package transaction

import (
	"olshop/product"
	"time"
)

type ServiceTrans struct {
	repo        RepoTransaction
	repoProduct product.RepoProduct
}

type ServiceTransInt interface {
	CreateTransaction(customerID, cartID int) (Transactions, error)
}

func NewTransactionService(repo RepoTransaction, repoProduct product.RepoProduct) *ServiceTrans {
	return &ServiceTrans{repo: repo, repoProduct: repoProduct}
}

func (s *ServiceTrans) CreateTransaction(customerID, cartID int) (Transactions, error) {
	products, err := s.repoProduct.GetListCartByID(cartID, customerID)
	if err != nil {
		return Transactions{}, err
	}

	idTrans, err := s.repo.GetLastTransaction()
	if err != nil {
		return Transactions{}, err
	}

	totalPrice := int32(0)
	for _, v := range products {
		totalPrice += v.Price
	}

	t := Transactions{
		ID:         idTrans,
		CustomerID: customerID,
		Price:      totalPrice,
		CreatedAt:  time.Now(),
		MaxTime:    time.Now().Add(time.Hour * 5),
		ShopCartID: cartID,
		PaymentID:  1,
	}

	err = s.repo.InserTransaction(t)
	if err != nil {
		return Transactions{}, err
	}

	createdTransaction, err := s.repo.GetDetailTransaction(idTrans)
	if err != nil {
		return Transactions{}, err
	}

	return createdTransaction, nil

}

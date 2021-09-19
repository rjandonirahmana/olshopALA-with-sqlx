package transaction

import (
	"fmt"
	"olshop/customer"
	"olshop/product"
	"time"
)

type ServiceTrans struct {
	repo        RepoTransaction
	repoProduct product.RepoProduct
}

type ServiceTransInt interface {
	CreateTransaction(customer customer.Customer, cartID int) (Transactions, error)
}

func NewTransactionService(repo RepoTransaction, repoProduct product.RepoProduct) *ServiceTrans {
	return &ServiceTrans{repo: repo, repoProduct: repoProduct}
}

func (s *ServiceTrans) CreateTransaction(customer customer.Customer, cartID int) (Transactions, error) {

	_, err := s.repoProduct.GetShopCartIDCustomer(customer.ID, cartID)
	if err != nil {
		return Transactions{}, err
	}

	products, err := s.repoProduct.GetListCartByID(cartID)
	if err != nil {
		return Transactions{}, err
	}

	id, err := s.repo.CheckTransaction(cartID)
	if err != nil {
		return Transactions{}, err
	}

	if id > 0 {
		transaction, err := s.repo.GetDetailTransaction(id)
		if err != nil {
			fmt.Println("error4")
			return Transactions{}, err
		}
		return transaction, nil
	}

	id, err = s.repo.GetLastIdTransaction()
	if err != nil {
		return Transactions{}, err
	}

	totalPrice := int32(0)
	for _, v := range products {
		totalPrice += v.Price * int32(v.Quantity)
	}

	t := Transactions{
		ID:         id + 1,
		CustomerID: customer.ID,
		Price:      totalPrice,
		CreatedAt:  time.Now(),
		MaxTime:    time.Now().Add(time.Hour * 5),
		ShopCartID: cartID,
		PaymentID:  1,
	}

	err = s.repo.InserTransaction(t)
	if err != nil {
		fmt.Println("error5")
		return Transactions{}, err
	}

	createdTransaction, err := s.repo.GetDetailTransaction(id + 1)
	if err != nil {
		fmt.Println("error6")
		return Transactions{}, err
	}

	return createdTransaction, nil

}

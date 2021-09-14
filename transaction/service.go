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
		fmt.Println("error1")
		return Transactions{}, err
	}

	products, err := s.repoProduct.GetListCartByID(cartID)
	if err != nil {
		fmt.Println("error2")
		return Transactions{}, err
	}

	checkid, _ := s.repo.CheckTransaction(cartID)

	if checkid > 0 {
		transaction, err := s.repo.GetDetailTransaction(checkid)
		if err != nil {
			fmt.Println("error4")
			return Transactions{}, err
		}
		return transaction, nil
	}

	totalPrice := int32(0)
	for _, v := range products {
		totalPrice += v.Price * int32(v.Quantity)
	}

	t := Transactions{
		ID:         0,
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

	createdTransaction, err := s.repo.GetDetailTransaction(checkid)
	if err != nil {
		fmt.Println("error6")
		return Transactions{}, err
	}

	return createdTransaction, nil

}

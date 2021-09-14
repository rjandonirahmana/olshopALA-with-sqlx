package transaction

import "time"

type Transactions struct {
	ID         int       `db:"id" json:"id"`
	CustomerID int       `db:"customer_id" json:"customer_id"`
	Price      int32     `db:"price" json:"price"`
	CreatedAt  time.Time `db:"created_at"`
	MaxTime    time.Time `db:"max_time" json:"time_limit"`
	ShopCartID int       `db:"shopcart_id" json:"shopcart_id"`
	PaymentID  int       `db:"payment_id" json:"payment_id"`
}

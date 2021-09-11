package transaction

import "time"

type Transactions struct {
	ID         int       `db:"id" json:"id"`
	ID_product int       `db:"product_id" json:"product_id"`
	CustomerID int       `db:"customer_id" json:"customer_id"`
	Quantity   int       `db:"quantity" json:"quantity"`
	Price      int32     `db:"price" json:"price"`
	CreatedAt  time.Time `db:"created_at"`
}

package shopcart

type Cart struct {
	CustomerID int `db:"customer_id" json:"customer_id"`
	ID         int `db:"id" json:"shopcart_id"`
}

type ShopeCart struct {
	CartID      int    `db:"cart_id" json:"cart_id"`
	ProductID   int    `db:"product_id" json:"product_id"`
	ProductName string `db:"product_name" json:"product_name"`
	Price       int32  `db:"price" json:"price"`
	Quantity    int    `db:"quantity" json:"quantity"`
}

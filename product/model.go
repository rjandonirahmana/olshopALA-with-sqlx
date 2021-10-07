package product

type Product struct {
	Name          string          `db:"name" json:"name"`
	ID            int64           `db:"id" json:"id"`
	Price         int32           `db:"price" json:"price"`
	Category_id   int             `db:"category_id" json:"category_id"`
	Description   *string         `db:"description" json:"description"`
	ProductImages ProductImage    `db:"product_images" json:"product_images,omitempty"`
	Quantity      int             `db:"quantity" json:"quantity"`
	SellerID      int             `db:"seller_id" json:"seller_id"`
	Category      ProductCategory `db:"product_category" json:"product_category"`
}

type ProductImage struct {
	ProductID *int    `db:"product_id" json:"product_id,omitempty"`
	IsPrimary *int    `db:"is_primary" json:"is_primary,omitempty"`
	Name      *string `db:"name" json:"name,omitempty"`
}

type ProductCategory struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Cart struct {
	CustomerID int `db:"customerID" json:"customer_id"`
	ID         int `db:"id" json:"shopcart_id"`
}

type ShopeCart struct {
	CartID      int    `db:"cart_id" json:"cart_id"`
	ProductID   int    `db:"product_id" json:"product_id"`
	ProductName string `db:"product_name" json:"product_name"`
	Price       int32  `db:"price" json:"price"`
	Quantity    int    `db:"quantity" json:"quantity"`
}

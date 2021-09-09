package product

type Product struct {
	Name          string       `db:"name" json:"name"`
	ID            int          `db:"id" json:"id"`
	Price         int32        `db:"price" json:"price"`
	Category_id   int          `db:"category_id" json:"category_id"`
	ProductImages ProductImage `db:"product_images" json:"product_images,omitempty"`
	Description   ProductDesc  `db:"product_desc" json:"description,omitempty"`
}

type ProductImage struct {
	ProductID int    `db:"product_id" json:"product_id"`
	IsPrimary int    `db:"is_primary" json:"is_primary"`
	Name      string `db:"name" json:"name"`
}

type ProductDesc struct {
	ProductID   int    `db:"product_id" json:"product_id"`
	Description string `db:"desc" json:"description"`
}

type ProductCategory struct {
	ID       int       `db:"id" json:"id"`
	Name     string    `db:"name" json:"name"`
	Products []Product `db:"products" json:"products"`
}

type Cart struct {
	CustomerID int `db:"customerID" json:"customer_id"`
	ID         int `db:"id" json:"id"`
}

type ShopeCart struct {
	CartID      int    `db:"cart_id" json:"cart_id"`
	ProductID   int    `db:"product_id" json:"product_id"`
	ProductName string `db:"product_name" json:"product_name"`
	Price       int32  `db:"price" json:"price"`
}

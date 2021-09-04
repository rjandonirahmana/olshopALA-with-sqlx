package product

type Product struct {
	Name          string       `db:"name" json:"name"`
	ID            int          `db:"id" json:"id"`
	Price         int          `db:"price" json:"price"`
	ProductImages ProductImage `db:"product_images" json:"product_images"`
	Description   ProductDesc  `db:"product_desc"`
}

type ProductImage struct {
	ProductID int    `db:"product_id" json:"product_id"`
	IsPrimary int    `db:"is_primary" json:"is_primary"`
	Name      string `db:"name" json:"name"`
}

type ProductDesc struct {
	ProductID   int     `db:"product_id" json:"product_id"`
	Description *string `db:"desc" json:"description"`
}

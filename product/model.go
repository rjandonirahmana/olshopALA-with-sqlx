package product

import "olshop/seller"

type Product struct {
	ID            uint            `db:"id" json:"id"`
	Name          string          `db:"name" json:"name"`
	Price         uint32          `db:"price" json:"price"`
	Quantity      uint            `db:"quantity" json:"quantity"`
	Description   *string         `db:"description" json:"description"`
	Rating        float32         `json:"rating"`
	SellerID      uint            `db:"seller_id" json:"-"`
	CategoryID    uint            `db:"category_id" json:"-"`
	Category      ProductCategory `db:"product_category" json:"category,omitempty"`
	Seller        seller.Seller   `db:"seller" json:"seller,omitempty"`
	ProductImages ProductImage    `db:"product_images" json:"product_images,omitempty"`
}

type ProductImage struct {
	ProductID uint   `db:"product_id" json:"product_id,omitempty"`
	IsPrimary uint   `db:"is_primary" json:"is_primary,omitempty"`
	Name      string `db:"name" json:"name,omitempty"`
}

type ProductCategory struct {
	ID   uint   `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}

type Products struct {
	Name          string          `db:"name" json:"name"`
	ID            uint            `db:"id" json:"id"`
	Price         int32           `db:"price" json:"price"`
	Quantity      uint            `db:"quantity" json:"quantity"`
	Category_id   int             `db:"category_id" json:"-"`
	Category      ProductCategory `db:"product_category" json:"category"`
	Description   *string         `db:"description" json:"description"`
	SellerID      uint            `db:"seller_id" json:"-"`
	ProductImages []ProductImage  `db:"product_images" json:"product_images,omitempty"`
}

package product

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type repoProduct struct {
	db *sqlx.DB
}

type RepoProduct interface {
	GetDetailsProductByID(id int) []Product
	GetProductByID(id int) (Product, error)
	GetProductByCategoryName(name_category string) ([]Product, error)
	InsertShoppingCart(cartid, productid int, price int32, name string) error
	GetLastID() (int, error)
	GetListCartByID(cartid int) ([]ShopeCart, error)
	CreateCart(customerID, id int) error
	GetShopCartIDCustomer(customerID, shopcartID int) (Cart, error)
	DeleteProductInShopCart(cart_id, customer_id, product_id int) error
	DecreaseQuantitInShopCart(cart_id, product_id int) error
	IncreaseQuantityInshopCart(cart_id, product_id int) error
	CheckInshopCart(cartid int, product_name string) (int, error)
	ShopCartCustomer(customerid int) ([]Cart, error)
	DeleteAllWhenQuantity0() error
}

func NewRepoProduct(db *sqlx.DB) *repoProduct {
	return &repoProduct{db: db}
}

func (r *repoProduct) GetDetailsProductByID(id int) []Product {
	querry := `SELECT p.id, p.name, p.description, p.quantity,pi.is_primary as "product_images.is_primary", pi.name as "product_images.name", pd.desc as "product_desc.desc" FROM products p INNER JOIN product_images pi ON p.id = pi.product_id INNER JOIN product_desc pd ON p.id = pd.product_id WHERE p.id = ? `

	var products []Product

	if err := r.db.Select(&products, querry, id); err != nil {
		fmt.Println(err)
		return []Product{}
	}

	return products
}

func (r *repoProduct) GetProductByID(id int) (Product, error) {
	querry := `SELECT * FROM products WHERE id = ?`

	var product Product
	err := r.db.Get(&product, querry, id)
	if err != nil {
		fmt.Println(err)
		return Product{}, err
	}

	return product, nil
}

func (r *repoProduct) GetProductByCategoryName(name_category string) ([]Product, error) {
	querry := `SELECT DISTINCT p.*, pi.name as "product_images.name", pi.is_primary as "product_images.is_primary", pi.product_id as "product_images.product_id", pc.id as "product_category.id", pc.name as "product_category.name" FROM products p LEFT JOIN product_category pc ON p.category_id = pc.id LEFT JOIN product_images pi ON p.id = pi.product_id WHERE pc.name = ? GROUP BY p.id`

	products := []Product{}

	err := r.db.Select(&products, querry, name_category)
	if err != nil {
		return []Product{}, err
	}

	return products, nil

}

func (r *repoProduct) CreateCart(customerID, id int) error {

	querry := `INSERT INTO cart (customerID, id) VALUES (?, ?)`

	_, err := r.db.Exec(querry, customerID, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repoProduct) ShopCartCustomer(customerid int) ([]Cart, error) {
	querry := `SELECT * FROM cart WHERE customerID = ?`

	var chart []Cart

	err := r.db.Select(&chart, querry, customerid)
	if err != nil {
		return []Cart{}, err
	}

	return chart, nil
}

func (r *repoProduct) InsertShoppingCart(cartid, productid int, price int32, name string) error {
	querry := `INSERT INTO shopcart (cart_id, product_id, product_name, price, quantity) VALUES(?, ?, ?, ?, ?)`

	_, err := r.db.Exec(querry, cartid, productid, name, price, 1)
	if err != nil {
		return err
	}
	return nil

}

func (r *repoProduct) GetLastID() (int, error) {
	querry := `SELECT id FROM cart WHERE id = (SELECT MAX(id) FROM cart)`

	var value int
	err := r.db.Get(&value, querry)
	if err != sql.ErrNoRows {
		return 0, err
	}
	return value, nil

}

func (r *repoProduct) GetListCartByID(cartid int) ([]ShopeCart, error) {
	querry := `SELECT cart_id, product_id, product_name, price, quantity FROM shopcart JOIN cart ON shopcart.cart_id = cart.id WHERE shopcart.cart_id = ?`

	var shopcart []ShopeCart
	err := r.db.Select(&shopcart, querry, cartid)
	if err != nil {
		return []ShopeCart{}, err
	}

	return shopcart, nil
}

func (r *repoProduct) GetShopCartIDCustomer(customerID, shopcartID int) (Cart, error) {

	querry := `SELECT * FROM cart WHERE customerID = ? AND id = ? ORDER BY id`

	var cart Cart

	err := r.db.Get(&cart, querry, customerID, shopcartID)

	if err != nil {
		return Cart{}, errors.New("customer doesnt have this cart id")
	}

	return cart, nil
}

func (r *repoProduct) DeleteProductInShopCart(cart_id, customer_id, product_id int) error {
	querry := `DELETE shopcart FROM shopcart JOIN cart ON shopcart.cart_id = cart.id WHERE shopcart.cart_id = ? AND cart.customerID = ? AND shopcart.product_id = ?`

	_, err := r.db.Exec(querry, cart_id, customer_id, product_id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repoProduct) DecreaseQuantitInShopCart(cart_id, product_id int) error {
	querry := `UPDATE shopcart SET quantity = quantity - 1 WHERE cart_id = ? AND product_id = ?`

	_, err := r.db.Exec(querry, cart_id, product_id)

	if err != nil {
		return err
	}
	return nil
}

func (r *repoProduct) IncreaseQuantityInshopCart(cart_id, product_id int) error {
	querry := `UPDATE shopcart SET quantity = quantity + 1 WHERE cart_id = ? AND product_id = ?`

	_, err := r.db.Exec(querry, cart_id, product_id)

	if err != nil {
		return err
	}
	return nil
}

func (r *repoProduct) CheckInshopCart(cartid int, product_name string) (int, error) {
	querry := `SELECT quantity FROM shopcart WHERE cart_id = ? AND product_name = ?`

	var quantity int
	err := r.db.Get(&quantity, querry, cartid, product_name)
	if err != nil {
		return 0, err
	}

	return quantity, nil
}

func (r *repoProduct) DeleteAllWhenQuantity0() error {
	querry := `DELETE FROM shopcart WHERE quantity = ?`

	_, err := r.db.Exec(querry, 0)

	if err != nil {
		return err
	}
	return nil
}

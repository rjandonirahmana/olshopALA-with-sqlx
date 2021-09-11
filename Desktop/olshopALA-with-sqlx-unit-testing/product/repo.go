package product

import (
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
	GetListCartByID(cartid int, customerid int) ([]Product, error)
	CreateCart(customerID, id int) error
	GetShopCartIDCustomer(customerID, shopcartID int) (Cart, error)
}

func NewRepoProduct(db *sqlx.DB) *repoProduct {
	return &repoProduct{db: db}
}

func (r *repoProduct) GetDetailsProductByID(id int) []Product {
	querry := `SELECT p.id, pi.is_primary as "product_images.is_primary", pi.name as "product_images.name", pd.desc as "product_desc.desc" FROM products p INNER JOIN product_images pi ON p.id = pi.product_id INNER JOIN product_desc pd ON p.id = pd.product_id WHERE p.id = ? `

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
	querry := `SELECT p.*, pi.name as "product_images.name", pi.is_primary as "product_images.is_primary", pi.product_id as "product_images.product_id" FROM products p LEFT JOIN product_category pc ON p.category_id = pc.id LEFT JOIN product_images pi ON pi.product_id = p.id WHERE pc.name = ?`

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

func (r *repoProduct) InsertShoppingCart(cartid, productid int, price int32, name string) error {
	querry := `INSERT INTO shopcart (cart_id, product_id, product_name, price) VALUES(?, ?, ?, ?)`

	_, err := r.db.Exec(querry, cartid, productid, name, price)
	if err != nil {
		return err
	}
	return nil

}

func (r *repoProduct) GetLastID() (int, error) {
	querry := `SELECT id FROM cart WHERE id = (SELECT MAX(id) FROM cart)`

	var value int
	err := r.db.Get(&value, querry)
	if err != nil {
		return 0, err
	}
	return value, nil

}

func (r *repoProduct) GetListCartByID(cartid int, customerid int) ([]Product, error) {
	querry := `SELECT p.name, p.price FROM products p JOIN shopcart sc ON p.id = sc.product_id JOIN cart c ON sc.cart_id = c.id WHERE sc.cart_id = ? AND c.customerID = ?`

	var products []Product
	err := r.db.Select(&products, querry, cartid, customerid)
	if err != nil {
		return []Product{}, err
	}

	return products, nil
}

func (r *repoProduct) DeleteProductInShopCart(cart_id, customer_id, product_id int) error {
	querry := `DELETE FROM shopcart JOIN cart ON shopcart.cart_id = cart.id WHERE shopcart.cart_id = ? AND cart.customerID = ? AND shopcart.product_id = ?`

	sqlx := r.db.QueryRowx(querry, cart_id, customer_id, product_id)
	if sqlx.Err() != nil {
		return sqlx.Err()
	}
	return nil
}

func (r *repoProduct) GetShopCartIDCustomer(customerID, shopcartID int) (Cart, error) {

	querry := `SELECT * FROM cart WHERE customerID = ? AND id = ?`

	var cart Cart

	err := r.db.Get(&cart, querry, customerID, shopcartID)

	if err != nil {
		fmt.Println(err)
		return Cart{}, errors.New(fmt.Sprintf("customer with id :%d doesnt have this cart with id : %d", customerID, shopcartID))
	}

	return cart, nil
}

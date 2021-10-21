package shopcart

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

type Repository interface {
	IncreaseQuantity(cart_id, product_id int) error
	InsertShoppingCart(cartid, productid int, price int32, name string) error
	ShopCartCustomer(customerid int) ([]Cart, error)
	GetListCartByID(cartid int) ([]ShopeCart, error)
	CreateCart(customerID int) error
	GetShopCartIDCustomer(customerID, shopcartID int) (Cart, error)
	DeleteProductInShopCart(cart_id, customer_id, product_id int) error
	DecreaseQuantitInShopCart(cart_id, product_id int) error
	CheckInshopCart(cartid int, product_name string) (int, error)
	DeleteAllWhenQuantity0() error
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) IncreaseQuantity(cart_id, product_id int) error {
	querry := `UPDATE shopcart SET quantity = quantity + 1 WHERE cart_id = $1 AND product_id = $2`

	_, err := r.db.Exec(querry, cart_id, product_id)

	if err != nil {
		return err
	}
	return nil
}

func (r *repository) InsertShoppingCart(cartid, productid int, price int32, name string) error {
	querry := `INSERT INTO shopcart (cart_id, product_id, product_name, price, quantity) VALUES(?, ?, ?, ?, ?)`

	_, err := r.db.Exec(querry, cartid, productid, name, price, 1)
	if err != nil {
		return err
	}
	return nil

}

func (r *repository) CreateCart(customerID int) error {

	querry := `INSERT INTO cart (customerID) VALUES (?)`

	_, err := r.db.Exec(querry, customerID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) ShopCartCustomer(customerid int) ([]Cart, error) {
	querry := `SELECT * FROM cart WHERE customerID = ?`

	var chart []Cart

	err := r.db.Select(&chart, querry, customerid)
	if err != nil {
		return []Cart{}, err
	}

	return chart, nil
}

func (r *repository) GetListCartByID(cartid int) ([]ShopeCart, error) {
	querry := `SELECT cart_id, product_id, product_name, price, quantity FROM shopcart JOIN cart ON shopcart.cart_id = cart.id WHERE shopcart.cart_id = $1`

	var shopcart []ShopeCart
	err := r.db.Select(&shopcart, querry, cartid)
	if err != nil {
		return []ShopeCart{}, err
	}

	return shopcart, nil
}

func (r *repository) GetShopCartIDCustomer(customerID, shopcartID int) (Cart, error) {

	querry := `SELECT * FROM cart WHERE customerID = ? AND id = ?`

	var cart Cart

	err := r.db.Get(&cart, querry, customerID, shopcartID)

	if err != sql.ErrNoRows {
		return cart, err
	}

	if cart.ID == 0 {
		return cart, errors.New("customer doesnt have this cart id")
	}

	return cart, nil
}

func (r *repository) DeleteProductInShopCart(cart_id, customer_id, product_id int) error {
	querry := `DELETE shopcart FROM shopcart JOIN cart ON shopcart.cart_id = cart.id WHERE shopcart.cart_id = ? AND cart.customerID = ? AND shopcart.product_id = ?`

	_, err := r.db.Exec(querry, cart_id, customer_id, product_id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DecreaseQuantitInShopCart(cart_id, product_id int) error {
	querry := `UPDATE shopcart SET quantity = quantity - 1 WHERE cart_id = ? AND product_id = ?`

	_, err := r.db.Exec(querry, cart_id, product_id)

	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CheckInshopCart(cartid int, product_name string) (int, error) {
	querry := `SELECT quantity FROM shopcart WHERE cart_id = ? AND product_name = ?`

	var quantity int
	err := r.db.Get(&quantity, querry, cartid, product_name)
	if err != nil {
		return 0, err
	}

	return quantity, nil
}

func (r *repository) DeleteAllWhenQuantity0() error {
	querry := `DELETE FROM shopcart WHERE quantity = ?`

	_, err := r.db.Exec(querry, 0)

	if err != nil {
		return err
	}
	return nil
}

package product

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type repoProduct struct {
	db *sqlx.DB
}

type RepoProduct interface {
	AddProduct(product Product) error
	InsertDetailProduct(product ProductDesc) error
	UpdatePrice(id int, price int) error
	GetProductsOrderByprice(name string) ([]Product, error)
	GetProductsOrderBypriceDSC(name string) ([]Product, error)
	GetDetailsProductByID(id int) []Product
	GetProductByID(id int) (Product, error)
	GetProductByCategoryName(name_category string) ([]Product, error)
	InsertShoppingCart(cartid, productid int, price int32, name string) error
	GetLastID() (int, error)
	GetListCartByID(cartid int, customerid int) ([]Product, error)
	CreateCart(customerID, id int) error
	DeleteProductInShopCart(cart_id, customer_id, product_id int) error
}

func NewRepoProduct(db *sqlx.DB) *repoProduct {
	return &repoProduct{db: db}
}

func (r *repoProduct) AddProduct(product Product) error {

	querry := `INSERT INTO products (id, name, price) VALUES (?, ?, ?) `

	_, err := r.db.Exec(querry, product.ID, product.Name, product.Price)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repoProduct) InsertDetailProduct(product ProductDesc) error {
	querry := `INSERT INTO product_desc (product_id, desc) VALUES (?, ?)`

	_, err := r.db.Exec(querry, product.ProductID, product.Description)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repoProduct) UpdatePrice(id int, price int) error {
	querry := `UPDATE products SET price = ? WHERE id = ?`

	_, err := r.db.Exec(querry, price, id)

	if err != nil {
		return err
	}

	return nil

}

func (r *repoProduct) GetProductsOrderByprice(name string) ([]Product, error) {
	querry := `SELECT * FROM products WHERE name = ? ORDER BY price`

	var products []Product

	err := r.db.Select(products, querry, name)
	if err != nil {
		return []Product{}, err
	}
	return products, nil

}

func (r *repoProduct) GetProductsOrderBypriceDSC(name string) ([]Product, error) {
	querry := `SELECT * FROM products WHERE name = ? ORDER BY price DESC`

	sql, err := r.db.Queryx(querry, name)

	if err != nil {
		fmt.Println(err)
		return []Product{}, err
	}

	var products []Product
	for sql.Next() {
		p := Product{}
		err := sql.StructScan(&p)

		if err != nil {
			fmt.Println(err)
			return []Product{}, err
		}

		products = append(products, p)
	}

	return products, nil
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
	querry := `SELECT p.name, p.price FROM products p INNER JOIN product_category pc ON p.category_id = pc.id WHERE pc.name = ?`

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
	querry := `DELETE shopcart FROM shopcart JOIN cart ON shopcart.cart_id = cart.id WHERE shopcart.cart_id = ? AND cart.customerID = ? AND shopcart.product_id = ?`

	_, err := r.db.Exec(querry, cart_id, customer_id, product_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repoProduct) SearchProducts(productname string) ([]Product, error) {
	querry := `SELECT  FROM products WHERE LOWER(name) LIKE ?`

	var products []Product
	err := r.db.Select(&products, querry, "%"+productname+"%")

	if err != nil {
		return []Product{}, err
	}
	return products, nil
}

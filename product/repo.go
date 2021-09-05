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
	GetDetailsProductByName(id int) []Product
	GetProductByID(id int) (Product, error)
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

func (r *repoProduct) GetDetailsProductByName(id int) []Product {
	querry := `SELECT p.id, pi.is_primary as "product_images.is_primary", pi.name as "product_images.name", pd.desc as "product_desc.desc" FROM products p INNER JOIN product_images pi ON p.id = pi.product_id INNER JOIN product_desc pd ON p.id = pd.product_id WHERE p.id = ? `

	var products []Product

	if err := r.db.Select(&products, querry, id); err != nil {
		fmt.Println(err)
		return []Product{}
	}

	return products
}

func (r *repoProduct) GetProductByID(id int) (Product, error) {
	querry := `SELECT * FROM product WHERE id = ?`

	var product Product
	err := r.db.Get(&product, querry, id)
	if err != nil {
		fmt.Println(err)
		return Product{}, err
	}

	return product, nil
}

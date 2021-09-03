package repo

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type repoProduct struct {
	db *sqlx.DB
}

func NewRepoProduct(db *sqlx.DB) *repoProduct {
	return &repoProduct{db: db}
}

func (r *repoProduct) InsertProduct(product Product) error {

	querry := `INSERT INTO products (id, name, price) VALUES (?, ?, ?) `

	_, err := r.db.Exec(querry, product.ID, product.Name, product.Price)

	if err != nil {
		return errors.New("eroro memasukan")
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

func (r *repoProduct) GetProductOrderByprice() ([]Product, error) {
	querry := `SELECT * FROM products ORDER BY price`

	sql, err := r.db.Queryx(querry)

	if err != nil {
		fmt.Println("disni db nya gabisa1")
		return []Product{}, err
	}

	var products []Product
	for sql.Next() {
		p := Product{}
		err := sql.StructScan(&p)

		if err != nil {
			fmt.Println("disni db nya gabisa2")
			return []Product{}, err
		}

		products = append(products, p)
	}

	return products, nil
}

package repository

import (
	"database/sql"

	"app/internal"
)

// NewProductsMySQL creates new mysql repository for product entity.
func NewProductsMySQL(db *sql.DB) *ProductsMySQL {
	return &ProductsMySQL{db}
}

// ProductsMySQL is the MySQL repository implementation for product entity.
type ProductsMySQL struct {
	// db is the database connection.
	db *sql.DB
}

// FindAll returns all products from the database.
func (r *ProductsMySQL) FindAll() (p []internal.Product, err error) {
	// execute the query
	rows, err := r.db.Query("SELECT `id`, `description`, `price` FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var pr internal.Product
		// scan the row into the product
		err := rows.Scan(&pr.Id, &pr.Description, &pr.Price)
		if err != nil {
			return nil, err
		}
		// append the product to the slice
		p = append(p, pr)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// Save saves the product into the database.
func (r *ProductsMySQL) Save(p *internal.Product) (err error) {
	// execute the query
	res, err := r.db.Exec(
		"INSERT INTO products (`description`, `price`) VALUES (?, ?)",
		(*p).Description, (*p).Price,
	)
	if err != nil {
		return err
	}

	// get the last inserted id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set the id
	(*p).Id = int(id)

	return
}

func (r *ProductsMySQL) FindTopProductsByAmount(limit int) ([]internal.ProductAmount, error) {
	var productsAmount []internal.ProductAmount
	rows, err := r.db.Query(
		"SELECT p.`description`, SUM(s.`quantity`) AS `total` "+
			"FROM products as p INNER JOIN sales as s ON p.`id` = s.`product_id` "+
			"GROUP BY p.`id` ORDER BY `total` DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pr internal.ProductAmount
		err := rows.Scan(&pr.Description, &pr.Total)
		if err != nil {
			return nil, err
		}
		productsAmount = append(productsAmount, pr)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return productsAmount, nil
}

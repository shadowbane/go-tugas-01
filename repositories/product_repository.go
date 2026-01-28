package repositories

import (
	"database/sql"
	"errors"

	"github.com/shadowbane/go-tugas-01/models"
)

var ErrProductNotFound = errors.New("product not found")

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	query := "SELECT id, product_id, name, description, price, stock FROM products"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.ProductId, &p.Name, &p.Description, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) GetByID(id string) (models.Product, error) {
	query := "SELECT id, product_id, name, description, price, stock FROM products WHERE id = $1"

	var p models.Product
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.ProductId, &p.Name, &p.Description, &p.Price, &p.Stock)
	if err == sql.ErrNoRows {
		return models.Product{}, ErrProductNotFound
	}
	if err != nil {
		return models.Product{}, err
	}

	return p, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (id, product_id, name, description, price, stock) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.db.Exec(query, product.ID, product.ProductId, product.Name, product.Description, product.Price, product.Stock)
	return err
}

func (r *ProductRepository) Update(id string, product *models.Product) error {
	query := "UPDATE products SET product_id = $1, name = $2, description = $3, price = $4, stock = $5 WHERE id = $6"
	result, err := r.db.Exec(query, product.ProductId, product.Name, product.Description, product.Price, product.Stock, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrProductNotFound
	}

	product.ID = id
	return nil
}

func (r *ProductRepository) Delete(id string) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrProductNotFound
	}

	return nil
}

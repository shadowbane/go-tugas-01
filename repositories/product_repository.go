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

func (r *ProductRepository) GetAll(name string) ([]models.Product, error) {
	query := "SELECT id, category_id, name, (SELECT name FROM categories WHERE categories.id = category_id) as category, description, price, stock FROM products"

	var args []interface{}
	if name != "" {
		query += " WHERE name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.CategoryId, &p.Name, &p.Category, &p.Description, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) GetByID(id string) (models.Product, error) {
	query := "SELECT id, category_id, name, (SELECT name FROM categories WHERE categories.id = category_id) as category, description, price, stock FROM products WHERE id = $1"

	var p models.Product
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.CategoryId, &p.Name, &p.Category, &p.Description, &p.Price, &p.Stock)
	if err == sql.ErrNoRows {
		return models.Product{}, ErrProductNotFound
	}
	if err != nil {
		return models.Product{}, err
	}

	return p, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "INSERT INTO products (id, category_id, name, description, price, stock) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = tx.Exec(query, product.ID, product.CategoryId, product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ProductRepository) Update(id string, product *models.Product) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lock the row first to prevent concurrent updates
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1 FOR UPDATE)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrProductNotFound
	}

	// Now safe to update - row is locked
	query := "UPDATE products SET category_id = $1, name = $2, description = $3, price = $4, stock = $5 WHERE id = $6"
	_, err = tx.Exec(query, product.CategoryId, product.Name, product.Description, product.Price, product.Stock, id)
	if err != nil {
		return err
	}

	product.ID = id
	return tx.Commit()
}

func (r *ProductRepository) Delete(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lock the row first to prevent concurrent modifications
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1 FOR UPDATE)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrProductNotFound
	}

	// Now safe to delete - row is locked
	_, err = tx.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

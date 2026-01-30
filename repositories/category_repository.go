package repositories

import (
	"database/sql"
	"errors"

	"github.com/shadowbane/go-tugas-01/models"
)

var ErrCategoryNotFound = errors.New("category not found")

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name, description FROM categories"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *CategoryRepository) GetByID(id string) (models.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"

	var c models.Category
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description)
	if err == sql.ErrNoRows {
		return models.Category{}, ErrCategoryNotFound
	}
	if err != nil {
		return models.Category{}, err
	}

	return c, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)"
	_, err = tx.Exec(query, category.ID, category.Name, category.Description)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CategoryRepository) Update(id string, category *models.Category) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lock the row first to prevent concurrent updates
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1 FOR UPDATE)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCategoryNotFound
	}

	// Now safe to update - row is locked
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	_, err = tx.Exec(query, category.Name, category.Description, id)
	if err != nil {
		return err
	}

	category.ID = id
	return tx.Commit()
}

func (r *CategoryRepository) Delete(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lock the row first to prevent concurrent modifications
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1 FOR UPDATE)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCategoryNotFound
	}

	// Now safe to delete - row is locked
	_, err = tx.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

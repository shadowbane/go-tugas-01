package repositories

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/shadowbane/go-tugas-01/dto"
	"github.com/shadowbane/go-tugas-01/models"
	"github.com/shadowbane/go-tugas-01/pkg/helpers"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrEmptyCart         = errors.New("cart is empty")
	ErrInvalidQuantity   = errors.New("quantity must be greater than 0")
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Checkout(items []dto.CheckoutItem) (*models.Transaction, error) {
	if len(items) == 0 {
		return nil, ErrEmptyCart
	}

	// Validate quantities
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var totalAmount int
	details := make([]models.TransactionDetail, 0, len(items))

	// Lock and validate each product, calculate totals
	for _, item := range items {
		var productID, productName string
		var price, stock int

		// Lock the row with SELECT FOR UPDATE
		err := tx.QueryRow(
			"SELECT id, name, price, stock FROM products WHERE id = $1 FOR UPDATE",
			item.ProductID,
		).Scan(&productID, &productName, &price, &stock)

		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound
		}
		if err != nil {
			return nil, err
		}

		// Check stock availability
		if stock < item.Quantity {
			return nil, ErrInsufficientStock
		}

		// Calculate subtotal
		subtotal := price * item.Quantity
		totalAmount += subtotal

		// Prepare transaction detail
		details = append(details, models.TransactionDetail{
			ID:          strings.ToUpper(helpers.NewULID()),
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Create transaction
	transaction := &models.Transaction{
		ID:          strings.ToUpper(helpers.NewULID()),
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		Details:     details,
	}

	// Insert transaction
	_, err = tx.Exec(
		"INSERT INTO transactions (id, total_amount, created_at) VALUES ($1, $2, $3)",
		transaction.ID, transaction.TotalAmount, transaction.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Insert details and update stock
	for i := range details {
		details[i].TransactionID = transaction.ID

		// Insert transaction detail
		_, err = tx.Exec(
			"INSERT INTO transaction_details (id, transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4, $5)",
			details[i].ID, details[i].TransactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal,
		)
		if err != nil {
			return nil, err
		}

		// Update product stock
		_, err = tx.Exec(
			"UPDATE products SET stock = stock - $1 WHERE id = $2",
			details[i].Quantity, details[i].ProductID,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction, nil
}

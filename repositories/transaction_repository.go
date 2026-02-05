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
	ErrInvalidDateRange  = errors.New("end_date must be after or equal to start_date")
)

// ReportResult holds the raw database result for report queries
type ReportResult struct {
	TotalRevenue   int
	TotalTransaksi int
	ProductName    sql.NullString
	TotalQty       sql.NullInt64
}

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

// GetTodayReport returns the report for today's transactions
func (r *TransactionRepository) GetTodayReport(includeDetails bool) (*ReportResult, []models.Transaction, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)
	return r.GetReportByDateRange(startOfDay, endOfDay, includeDetails)
}

// GetReportByDateRange returns transaction report for a given date range
func (r *TransactionRepository) GetReportByDateRange(startDate, endDate time.Time, includeDetails bool) (*ReportResult, []models.Transaction, error) {
	if endDate.Before(startDate) {
		return nil, nil, ErrInvalidDateRange
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	// CTE query for report summary
	query := `
		WITH date_transactions AS (
			SELECT t.id, t.total_amount, t.created_at
			FROM transactions t
			WHERE t.created_at >= $1 AND t.created_at < $2
		),
		transaction_summary AS (
			SELECT
				COALESCE(SUM(total_amount), 0) as total_revenue,
				COUNT(*) as total_transaksi
			FROM date_transactions
		),
		best_product AS (
			SELECT
				COALESCE(p.name, td.product_id) as product_name,
				SUM(td.quantity) as total_qty
			FROM transaction_details td
			JOIN date_transactions dt ON td.transaction_id = dt.id
			LEFT JOIN products p ON td.product_id = p.id
			GROUP BY td.product_id, p.name
			ORDER BY total_qty DESC
			LIMIT 1
		)
		SELECT ts.total_revenue, ts.total_transaksi, bp.product_name, bp.total_qty
		FROM transaction_summary ts
		LEFT JOIN best_product bp ON true
	`

	result := &ReportResult{}
	err = tx.QueryRow(query, startDate, endDate).Scan(
		&result.TotalRevenue,
		&result.TotalTransaksi,
		&result.ProductName,
		&result.TotalQty,
	)
	if err != nil {
		return nil, nil, err
	}

	var transactions []models.Transaction
	if includeDetails {
		transactions, err = r.getTransactionsInRange(tx, startDate, endDate)
		if err != nil {
			return nil, nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	return result, transactions, nil
}

// getTransactionsInRange retrieves all transactions with their details within a date range
func (r *TransactionRepository) getTransactionsInRange(tx *sql.Tx, startDate, endDate time.Time) ([]models.Transaction, error) {
	// Get transactions
	rows, err := tx.Query(
		"SELECT id, total_amount, created_at FROM transactions WHERE created_at >= $1 AND created_at < $2 ORDER BY created_at DESC",
		startDate, endDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.TotalAmount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Get details for each transaction
	for i := range transactions {
		detailRows, err := tx.Query(
			`SELECT td.id, td.transaction_id, td.product_id, COALESCE(p.name, td.product_id) as product_name, td.quantity, td.subtotal
			 FROM transaction_details td
			 LEFT JOIN products p ON td.product_id = p.id
			 WHERE td.transaction_id = $1`,
			transactions[i].ID,
		)
		if err != nil {
			return nil, err
		}

		var details []models.TransactionDetail
		for detailRows.Next() {
			var d models.TransactionDetail
			if err := detailRows.Scan(&d.ID, &d.TransactionID, &d.ProductID, &d.ProductName, &d.Quantity, &d.Subtotal); err != nil {
				detailRows.Close()
				return nil, err
			}
			details = append(details, d)
		}
		detailRows.Close()

		if err := detailRows.Err(); err != nil {
			return nil, err
		}

		transactions[i].Details = details
	}

	return transactions, nil
}

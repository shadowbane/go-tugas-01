package services

import (
	"sort"
	"time"

	"github.com/shadowbane/go-tugas-01/dto"
	"github.com/shadowbane/go-tugas-01/models"
	"github.com/shadowbane/go-tugas-01/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{
		repo: repo,
	}
}

func (s *TransactionService) Checkout(request *dto.CheckoutRequest) (*models.Transaction, error) {
	return s.repo.Checkout(request.Items)
}

// GetTodayReport returns today's report in consolidated format
func (s *TransactionService) GetTodayReport(includeDetails bool) (*dto.ReportResponse, error) {
	result, transactions, err := s.repo.GetTodayReport(includeDetails)
	if err != nil {
		return nil, err
	}
	return s.buildConsolidatedResponse(result, transactions, includeDetails), nil
}

// GetReportByDateRange returns report for date range
// consolidated=true: returns *dto.ReportResponse
// consolidated=false: returns []dto.DailyReport (sliced by date)
func (s *TransactionService) GetReportByDateRange(startDate, endDate time.Time, includeDetails, consolidated bool) (interface{}, error) {
	// Always fetch transactions (needed for daily aggregates in sliced mode)
	result, transactions, err := s.repo.GetReportByDateRange(startDate, endDate, true)
	if err != nil {
		return nil, err
	}

	if consolidated {
		return s.buildConsolidatedResponse(result, transactions, includeDetails), nil
	}
	return s.buildSlicedResponse(transactions, includeDetails), nil
}

// buildConsolidatedResponse builds a flat response with all transaction details
func (s *TransactionService) buildConsolidatedResponse(result *repositories.ReportResult, transactions []models.Transaction, includeDetails bool) *dto.ReportResponse {
	response := &dto.ReportResponse{
		TotalRevenue:   result.TotalRevenue,
		TotalTransaksi: result.TotalTransaksi,
	}

	if result.ProductName.Valid && result.TotalQty.Valid {
		response.ProdukTerlaris = &dto.ProdukTerlaris{
			Nama:       result.ProductName.String,
			QtyTerjual: int(result.TotalQty.Int64),
		}
	}

	if includeDetails {
		// Flatten all transaction details into one array
		var allDetails []models.TransactionDetail
		for _, t := range transactions {
			allDetails = append(allDetails, t.Details...)
		}
		response.Details = allDetails
	}

	return response
}

// buildSlicedResponse builds a response sliced by date with daily aggregates
func (s *TransactionService) buildSlicedResponse(transactions []models.Transaction, includeDetails bool) []dto.DailyReport {
	if len(transactions) == 0 {
		return []dto.DailyReport{}
	}

	// Group by date
	dateMap := make(map[string][]models.Transaction)
	for _, t := range transactions {
		dateKey := t.CreatedAt.Format("2006-01-02")
		dateMap[dateKey] = append(dateMap[dateKey], t)
	}

	// Convert map to slice with daily aggregates
	result := make([]dto.DailyReport, 0, len(dateMap))
	for date, txns := range dateMap {
		dailyReport := dto.DailyReport{
			Date:           date,
			TotalTransaksi: len(txns),
		}

		// Calculate daily revenue and find best-selling product
		productQty := make(map[string]int)     // productID -> total qty
		productName := make(map[string]string) // productID -> product name
		var allDetails []models.TransactionDetail

		for _, t := range txns {
			dailyReport.TotalRevenue += t.TotalAmount
			for _, d := range t.Details {
				productQty[d.ProductID] += d.Quantity
				productName[d.ProductID] = d.ProductName
				allDetails = append(allDetails, d)
			}
		}

		// Find best-selling product for this day
		var bestProductID string
		var maxQty int
		for pid, qty := range productQty {
			if qty > maxQty {
				maxQty = qty
				bestProductID = pid
			}
		}

		if bestProductID != "" {
			dailyReport.ProdukTerlaris = &dto.ProdukTerlaris{
				Nama:       productName[bestProductID],
				QtyTerjual: maxQty,
			}
		}

		// Include transaction details (line items) only if requested
		if includeDetails {
			dailyReport.Details = allDetails
		}

		result = append(result, dailyReport)
	}

	// Sort by date descending (newest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date > result[j].Date
	})

	return result
}

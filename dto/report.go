package dto

import "github.com/shadowbane/go-tugas-01/models"

type ProdukTerlaris struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

// DailyReport represents a daily summary with aggregates (for sliced mode)
type DailyReport struct {
	Date           string                     `json:"date"`
	TotalRevenue   int                        `json:"total_revenue"`
	TotalTransaksi int                        `json:"total_transaksi"`
	ProdukTerlaris *ProdukTerlaris            `json:"produk_terlaris"`
	Details        []models.TransactionDetail `json:"details,omitempty"`
}

// ReportResponse represents consolidated report format
type ReportResponse struct {
	TotalRevenue   int                        `json:"total_revenue"`
	TotalTransaksi int                        `json:"total_transaksi"`
	ProdukTerlaris *ProdukTerlaris            `json:"produk_terlaris"`
	Details        []models.TransactionDetail `json:"details,omitempty"`
}

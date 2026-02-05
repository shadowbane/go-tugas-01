package services

import (
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

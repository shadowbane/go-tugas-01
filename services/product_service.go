package services

import (
	"strings"

	"github.com/shadowbane/go-tugas-01/models"
	"github.com/shadowbane/go-tugas-01/pkg/helpers"
	"github.com/shadowbane/go-tugas-01/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetByID(id string) (models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Create(product *models.Product) error {
	product.ID = strings.ToUpper(helpers.NewULID())
	return s.repo.Create(product)
}

func (s *ProductService) Update(id string, product *models.Product) error {
	return s.repo.Update(id, product)
}

func (s *ProductService) Delete(id string) error {
	return s.repo.Delete(id)
}

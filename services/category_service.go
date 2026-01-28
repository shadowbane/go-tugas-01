package services

import (
	"strings"

	"github.com/shadowbane/go-tugas-01/models"
	"github.com/shadowbane/go-tugas-01/pkg/helpers"
	"github.com/shadowbane/go-tugas-01/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) GetByID(id string) (models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Create(category *models.Category) error {
	category.ID = strings.ToUpper(helpers.NewULID())
	return s.repo.Create(category)
}

func (s *CategoryService) Update(id string, category *models.Category) error {
	return s.repo.Update(id, category)
}

func (s *CategoryService) Delete(id string) error {
	return s.repo.Delete(id)
}

package services

import (
	"cashier-api/models"
	"cashier-api/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll(name string) ([]models.Category, error) {
	return s.repo.GetAll(name)
}

func (s *CategoryService) Create(category *models.Category) error {
	return s.repo.Create(category)
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Update(category *models.Category) (*models.Category, error) {
	err := s.repo.Update(category)
	if err != nil {
		return nil, err
	}
	return s.GetByID(category.ID)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}

package services

import (
	"cashier-api/models"
	"cashier-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.repo.GetAll(name)
}

func (s *ProductService) Create(data *models.Product) error {
	return s.repo.Create(data)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(product *models.Product) (*models.Product, error) {
	err := s.repo.Update(product)
	if err != nil {
		return nil, err
	}
	return s.GetByID(product.ID)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}

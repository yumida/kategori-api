package services

import (
	"cashier-api/models"
	"cashier-api/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetReportToday() (*models.ReportResponse, error) {
	return s.repo.GetReportToday()
}

func (s *TransactionService) GetReportByDate(StartDate string, EndDate string) (*models.ReportResponse, error) {
	return s.repo.GetReportByDate(StartDate, EndDate)
}

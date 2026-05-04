package service

import (
	"fmt"
	"exchangeapp/internal/model"
	"exchangeapp/internal/repository"
	"time"
)

type ExchangeRateService struct {
	repo repository.ExchangeRateRepository
}

func NewExchangeRateService(repo repository.ExchangeRateRepository) *ExchangeRateService {
	return &ExchangeRateService{repo: repo}
}

func (s *ExchangeRateService) CreateExchangeRate(rate *model.ExchangeRate) error {
	rate.Date = time.Now()
	if err := s.repo.Create(rate); err != nil {
		return fmt.Errorf("exchangeRepo.Create: %w", err)
	}
	return nil
}

func (s *ExchangeRateService) GetExchangeRates() ([]model.ExchangeRate, error) {
	rates, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("exchangeRepo.FindAll: %w", err)
	}
	return rates, nil
}

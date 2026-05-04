package service

import (
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
	return s.repo.Create(rate)
}

func (s *ExchangeRateService) GetExchangeRates() ([]model.ExchangeRate, error) {
	return s.repo.FindAll()
}

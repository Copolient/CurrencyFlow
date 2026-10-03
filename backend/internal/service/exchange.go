package service

import (
	"fmt"
	"strings"
	"time"

	"currencyflow/internal/model"
	"currencyflow/internal/repository"
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

// UpsertExchangeRate keeps the live snapshot table in sync with collected rates.
func (s *ExchangeRateService) UpsertExchangeRate(from, to string, rate float64) error {
	r := &model.ExchangeRate{
		FromCurrency: strings.ToUpper(strings.TrimSpace(from)),
		ToCurrency:   strings.ToUpper(strings.TrimSpace(to)),
		Rate:         rate,
		Date:         time.Now(),
	}
	if err := s.repo.Upsert(r); err != nil {
		return fmt.Errorf("exchangeRepo.Upsert: %w", err)
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

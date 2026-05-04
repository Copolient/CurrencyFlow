package service_test

import (
	"errors"
	"exchangeapp/internal/mock"
	"exchangeapp/internal/model"
	"exchangeapp/internal/service"
	"testing"
)

func TestCreateExchangeRate_Success(t *testing.T) {
	repo := mock.NewExchangeRateRepo()
	svc := service.NewExchangeRateService(repo)

	rate := &model.ExchangeRate{FromCurrency: "USD", ToCurrency: "CNY", Rate: 7.24}
	err := svc.CreateExchangeRate(rate)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rate.Date.IsZero() {
		t.Fatal("expected Date to be set")
	}
}

func TestCreateExchangeRate_RepoError(t *testing.T) {
	repo := mock.NewExchangeRateRepo()
	repo.Err = errors.New("write failed")
	svc := service.NewExchangeRateService(repo)

	err := svc.CreateExchangeRate(&model.ExchangeRate{FromCurrency: "USD", ToCurrency: "CNY", Rate: 7.24})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetExchangeRates_Success(t *testing.T) {
	repo := mock.NewExchangeRateRepo()
	svc := service.NewExchangeRateService(repo)

	_ = svc.CreateExchangeRate(&model.ExchangeRate{FromCurrency: "USD", ToCurrency: "CNY", Rate: 7.24})
	_ = svc.CreateExchangeRate(&model.ExchangeRate{FromCurrency: "EUR", ToCurrency: "USD", Rate: 1.08})

	rates, err := svc.GetExchangeRates()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(rates) != 2 {
		t.Fatalf("expected 2 rates, got %d", len(rates))
	}
}

func TestGetExchangeRates_Empty(t *testing.T) {
	repo := mock.NewExchangeRateRepo()
	svc := service.NewExchangeRateService(repo)

	rates, err := svc.GetExchangeRates()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(rates) != 0 {
		t.Fatalf("expected 0 rates, got %d", len(rates))
	}
}

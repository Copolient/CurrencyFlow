package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"exchangeapp/internal/mock"
	"exchangeapp/internal/model"
	"exchangeapp/internal/service"
)

func TestRecordRate_Success(t *testing.T) {
	repo := mock.NewRateHistoryRepo()
	cache := mock.NewCache()
	svc := service.NewRateHistoryService(repo, cache)

	err := svc.RecordRate(context.Background(), "USD", "CNY", 7.24)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRecordRate_RepoError(t *testing.T) {
	repo := mock.NewRateHistoryRepo()
	repo.Err = errors.New("db write failed")
	cache := mock.NewCache()
	svc := service.NewRateHistoryService(repo, cache)

	err := svc.RecordRate(context.Background(), "USD", "CNY", 7.24)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetHistoryByPair_CacheHit(t *testing.T) {
	repo := mock.NewRateHistoryRepo()
	cache := mock.NewCache()
	svc := service.NewRateHistoryService(repo, cache)

	// Pre-populate cache
	histories := []model.ExchangeRateHistory{
		{FromCurrency: "USD", ToCurrency: "CNY", Rate: 7.24, Timestamp: time.Now()},
	}
	data, _ := json.Marshal(histories)
	cache.Set(context.Background(), "rate_history:USD:CNY:1M", string(data), 5*time.Minute)

	result, err := svc.GetHistoryByPair(context.Background(), "USD", "CNY", "1M")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 history, got %d", len(result))
	}
	if result[0].Rate != 7.24 {
		t.Fatalf("expected rate 7.24, got %f", result[0].Rate)
	}
}

func TestGetHistoryByPair_CacheMiss(t *testing.T) {
	repo := mock.NewRateHistoryRepo()
	cache := mock.NewCache()
	svc := service.NewRateHistoryService(repo, cache)

	// Record some rates
	now := time.Now()
	repo.Create(&model.ExchangeRateHistory{FromCurrency: "USD", ToCurrency: "CNY", Rate: 7.20, Timestamp: now.Add(-2 * 24 * time.Hour)})
	repo.Create(&model.ExchangeRateHistory{FromCurrency: "USD", ToCurrency: "CNY", Rate: 7.24, Timestamp: now})

	result, err := svc.GetHistoryByPair(context.Background(), "USD", "CNY", "1M")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 histories, got %d", len(result))
	}
}

func TestGetHistoryByPair_Empty(t *testing.T) {
	repo := mock.NewRateHistoryRepo()
	cache := mock.NewCache()
	svc := service.NewRateHistoryService(repo, cache)

	result, err := svc.GetHistoryByPair(context.Background(), "USD", "CNY", "1M")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 histories, got %d", len(result))
	}
}

func TestGetLatestByAllPairs_Success(t *testing.T) {
	repo := mock.NewRateHistoryRepo()
	cache := mock.NewCache()
	svc := service.NewRateHistoryService(repo, cache)

	now := time.Now()
	repo.Create(&model.ExchangeRateHistory{FromCurrency: "USD", ToCurrency: "CNY", Rate: 7.20, Timestamp: now.Add(-1 * time.Hour)})
	repo.Create(&model.ExchangeRateHistory{FromCurrency: "USD", ToCurrency: "CNY", Rate: 7.24, Timestamp: now})
	repo.Create(&model.ExchangeRateHistory{FromCurrency: "EUR", ToCurrency: "USD", Rate: 1.08, Timestamp: now})

	result, err := svc.GetLatestByAllPairs(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(result))
	}

	// Verify USD/CNY has the latest rate
	for _, h := range result {
		if h.FromCurrency == "USD" && h.ToCurrency == "CNY" {
			if h.Rate != 7.24 {
				t.Fatalf("expected latest USD/CNY rate 7.24, got %f", h.Rate)
			}
		}
	}
}

func TestGetLatestByAllPairs_CacheHit(t *testing.T) {
	repo := mock.NewRateHistoryRepo()
	cache := mock.NewCache()
	svc := service.NewRateHistoryService(repo, cache)

	// Pre-populate cache
	histories := []model.ExchangeRateHistory{
		{FromCurrency: "USD", ToCurrency: "CNY", Rate: 7.24, Timestamp: time.Now()},
	}
	data, _ := json.Marshal(histories)
	cache.Set(context.Background(), "rate_history:latest:all", string(data), 5*time.Minute)

	result, err := svc.GetLatestByAllPairs(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(result))
	}
}

func TestParseTimeRange(t *testing.T) {
	now := time.Now()

	tests := []struct {
		rangeStr string
		wantDays int
	}{
		{"1D", 1},
		{"1W", 7},
		{"1M", 30},
		{"3M", 90},
		{"1Y", 365},
		{"invalid", 30}, // default
	}

	for _, tt := range tests {
		t.Run(tt.rangeStr, func(t *testing.T) {
			repo := mock.NewRateHistoryRepo()
			cache := mock.NewCache()
			svc := service.NewRateHistoryService(repo, cache)

			// Create history within the expected range
			historyTime := now.Add(-time.Duration(tt.wantDays/2) * 24 * time.Hour)
			repo.Create(&model.ExchangeRateHistory{
				FromCurrency: "USD", ToCurrency: "CNY", Rate: 7.24, Timestamp: historyTime,
			})

			result, err := svc.GetHistoryByPair(context.Background(), "USD", "CNY", tt.rangeStr)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			// The history should be within range
			if len(result) == 0 {
				t.Fatalf("expected history to be within range %s", tt.rangeStr)
			}
		})
	}
}

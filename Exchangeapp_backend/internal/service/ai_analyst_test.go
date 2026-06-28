package service_test

import (
	"context"
	"testing"
	"time"

	"exchangeapp/internal/mock"
	"exchangeapp/internal/model"
	"exchangeapp/internal/service"
	"exchangeapp/pkg/config"
)

func TestAnalyze_WithData(t *testing.T) {
	repo := mock.NewRateHistoryRepo()
	cache := mock.NewCache()
	svc := service.NewAIAnalystService(repo, cache, config.LLMConfig{})

	// Create some history data
	now := time.Now()
	for i := 0; i < 30; i++ {
		rate := 7.0 + float64(i)*0.01 // upward trend
		repo.Create(&model.ExchangeRateHistory{
			FromCurrency: "USD",
			ToCurrency:   "CNY",
			Rate:         rate,
			Timestamp:    now.Add(-time.Duration(30-i) * 24 * time.Hour),
		})
	}

	result, err := svc.Analyze(context.Background(), "USD", "CNY", "美元还会涨吗?")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Analysis == "" {
		t.Fatal("expected analysis to be non-empty")
	}
	if result.Trend != "bullish" {
		t.Fatalf("expected bullish trend, got %s", result.Trend)
	}
	if result.RiskWarning == "" {
		t.Fatal("expected risk warning to be non-empty")
	}
	if result.KeyLevels.Support >= result.KeyLevels.Resistance {
		t.Fatal("expected support < resistance")
	}
}

func TestAnalyze_NoData(t *testing.T) {
	repo := mock.NewRateHistoryRepo()
	cache := mock.NewCache()
	svc := service.NewAIAnalystService(repo, cache, config.LLMConfig{})

	result, err := svc.Analyze(context.Background(), "USD", "CNY", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Trend != "neutral" {
		t.Fatalf("expected neutral trend for no data, got %s", result.Trend)
	}
}

func TestAnalyze_BearishTrend(t *testing.T) {
	repo := mock.NewRateHistoryRepo()
	cache := mock.NewCache()
	svc := service.NewAIAnalystService(repo, cache, config.LLMConfig{})

	// Create downward trend
	now := time.Now()
	for i := 0; i < 30; i++ {
		rate := 7.5 - float64(i)*0.01 // downward trend
		repo.Create(&model.ExchangeRateHistory{
			FromCurrency: "USD",
			ToCurrency:   "CNY",
			Rate:         rate,
			Timestamp:    now.Add(-time.Duration(30-i) * 24 * time.Hour),
		})
	}

	result, err := svc.Analyze(context.Background(), "USD", "CNY", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Trend != "bearish" {
		t.Fatalf("expected bearish trend, got %s", result.Trend)
	}
}

func TestAnalyze_CacheHit(t *testing.T) {
	repo := mock.NewRateHistoryRepo()
	cache := mock.NewCache()
	svc := service.NewAIAnalystService(repo, cache, config.LLMConfig{})

	// Create data and call once to populate cache
	now := time.Now()
	for i := 0; i < 10; i++ {
		repo.Create(&model.ExchangeRateHistory{
			FromCurrency: "USD",
			ToCurrency:   "CNY",
			Rate:         7.0 + float64(i)*0.01,
			Timestamp:    now.Add(-time.Duration(10-i) * 24 * time.Hour),
		})
	}

	// First call populates cache
	result1, err := svc.Analyze(context.Background(), "USD", "CNY", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Second call should hit cache
	result2, err := svc.Analyze(context.Background(), "USD", "CNY", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result2.Analysis != result1.Analysis {
		t.Fatalf("expected cached analysis to match, got different results")
	}
}

func TestAnalyze_WithQuestion(t *testing.T) {
	repo := mock.NewRateHistoryRepo()
	cache := mock.NewCache()
	svc := service.NewAIAnalystService(repo, cache, config.LLMConfig{})

	// Create some data
	now := time.Now()
	repo.Create(&model.ExchangeRateHistory{
		FromCurrency: "EUR",
		ToCurrency:   "USD",
		Rate:         1.08,
		Timestamp:    now.Add(-24 * time.Hour),
	})
	repo.Create(&model.ExchangeRateHistory{
		FromCurrency: "EUR",
		ToCurrency:   "USD",
		Rate:         1.09,
		Timestamp:    now,
	})

	result, err := svc.Analyze(context.Background(), "EUR", "USD", "欧元走势如何?")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Analysis == "" {
		t.Fatal("expected analysis to be non-empty")
	}
}

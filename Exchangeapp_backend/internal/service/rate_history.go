package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"exchangeapp/internal/model"
	"exchangeapp/internal/repository"
	"exchangeapp/pkg/cache"
)

type RateHistoryService struct {
	repo  repository.RateHistoryRepository
	cache cache.Cache
}

func NewRateHistoryService(repo repository.RateHistoryRepository, cache cache.Cache) *RateHistoryService {
	return &RateHistoryService{repo: repo, cache: cache}
}

const rateHistoryCacheTTL = 5 * time.Minute

func (s *RateHistoryService) GetHistoryByPair(ctx context.Context, from, to string, rangeStr string) ([]model.ExchangeRateHistory, error) {
	cacheKey := fmt.Sprintf("rate_history:%s:%s:%s", from, to, rangeStr)

	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var histories []model.ExchangeRateHistory
			if json.Unmarshal([]byte(cached), &histories) == nil {
				return histories, nil
			}
		}
	}

	end := time.Now()
	start := parseTimeRange(rangeStr, end)

	histories, err := s.repo.FindByPairAndTimeRange(from, to, start, end)
	if err != nil {
		return nil, fmt.Errorf("rateHistoryRepo.FindByPairAndTimeRange: %w", err)
	}

	if s.cache != nil && len(histories) > 0 {
		data, _ := json.Marshal(histories)
		_ = s.cache.Set(ctx, cacheKey, string(data), rateHistoryCacheTTL)
	}

	return histories, nil
}

func (s *RateHistoryService) GetLatestByAllPairs(ctx context.Context) ([]model.ExchangeRateHistory, error) {
	cacheKey := "rate_history:latest:all"

	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var histories []model.ExchangeRateHistory
			if json.Unmarshal([]byte(cached), &histories) == nil {
				return histories, nil
			}
		}
	}

	histories, err := s.repo.FindLatestByAllPairs()
	if err != nil {
		return nil, fmt.Errorf("rateHistoryRepo.FindLatestByAllPairs: %w", err)
	}

	if s.cache != nil && len(histories) > 0 {
		data, _ := json.Marshal(histories)
		_ = s.cache.Set(ctx, cacheKey, string(data), rateHistoryCacheTTL)
	}

	return histories, nil
}

func (s *RateHistoryService) RecordRate(ctx context.Context, from, to string, rate float64) error {
	history := &model.ExchangeRateHistory{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		Timestamp:    time.Now(),
	}
	if err := s.repo.Create(history); err != nil {
		return fmt.Errorf("rateHistoryRepo.Create: %w", err)
	}

	// Invalidate cache
	if s.cache != nil {
		_ = s.cache.Del(ctx, "rate_history:latest:all")
	}

	return nil
}

func parseTimeRange(rangeStr string, end time.Time) time.Time {
	switch rangeStr {
	case "1D":
		return end.Add(-24 * time.Hour)
	case "1W":
		return end.Add(-7 * 24 * time.Hour)
	case "1M":
		return end.Add(-30 * 24 * time.Hour)
	case "3M":
		return end.Add(-90 * 24 * time.Hour)
	case "1Y":
		return end.Add(-365 * 24 * time.Hour)
	default:
		return end.Add(-30 * 24 * time.Hour) // default 1 month
	}
}

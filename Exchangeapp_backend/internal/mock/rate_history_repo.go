package mock

import (
	"exchangeapp/internal/model"
	"sync"
	"time"
)

type RateHistoryRepo struct {
	mu        sync.RWMutex
	histories []model.ExchangeRateHistory
	Err       error
}

func NewRateHistoryRepo() *RateHistoryRepo {
	return &RateHistoryRepo{}
}

func (r *RateHistoryRepo) Create(history *model.ExchangeRateHistory) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.histories = append(r.histories, *history)
	return nil
}

func (r *RateHistoryRepo) BulkCreate(histories []model.ExchangeRateHistory) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.histories = append(r.histories, histories...)
	return nil
}

func (r *RateHistoryRepo) FindByPair(from, to string) ([]model.ExchangeRateHistory, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.ExchangeRateHistory
	for _, h := range r.histories {
		if h.FromCurrency == from && h.ToCurrency == to {
			result = append(result, h)
		}
	}
	return result, nil
}

func (r *RateHistoryRepo) FindByPairAndTimeRange(from, to string, start, end time.Time) ([]model.ExchangeRateHistory, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.ExchangeRateHistory
	for _, h := range r.histories {
		if h.FromCurrency == from && h.ToCurrency == to &&
			(h.Timestamp.After(start) || h.Timestamp.Equal(start)) &&
			(h.Timestamp.Before(end) || h.Timestamp.Equal(end)) {
			result = append(result, h)
		}
	}
	return result, nil
}

func (r *RateHistoryRepo) FindLatestByAllPairs() ([]model.ExchangeRateHistory, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Group by pair, keep latest
	latest := make(map[string]model.ExchangeRateHistory)
	for _, h := range r.histories {
		key := h.FromCurrency + ":" + h.ToCurrency
		if existing, ok := latest[key]; !ok || h.Timestamp.After(existing.Timestamp) {
			latest[key] = h
		}
	}
	var result []model.ExchangeRateHistory
	for _, h := range latest {
		result = append(result, h)
	}
	return result, nil
}

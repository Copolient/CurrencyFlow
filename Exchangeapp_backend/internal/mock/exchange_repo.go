package mock

import (
	"exchangeapp/internal/model"
	"sync"
)

type ExchangeRateRepo struct {
	mu    sync.RWMutex
	rates []model.ExchangeRate
	Err   error
}

func NewExchangeRateRepo() *ExchangeRateRepo {
	return &ExchangeRateRepo{}
}

func (r *ExchangeRateRepo) Create(rate *model.ExchangeRate) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rates = append(r.rates, *rate)
	return nil
}

func (r *ExchangeRateRepo) FindAll() ([]model.ExchangeRate, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.rates, nil
}

package mock

import (
	"currencyflow/internal/model"
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

func (r *ExchangeRateRepo) Upsert(rate *model.ExchangeRate) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.rates {
		if r.rates[i].FromCurrency == rate.FromCurrency && r.rates[i].ToCurrency == rate.ToCurrency {
			r.rates[i] = *rate
			return nil
		}
	}
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

package mock

import (
	"exchangeapp/internal/model"
	"sync"
)

type AlertRepo struct {
	mu     sync.RWMutex
	alerts []model.RateAlert
	Err    error
}

func NewAlertRepo() *AlertRepo {
	return &AlertRepo{}
}

func (r *AlertRepo) Create(alert *model.RateAlert) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	alert.ID = uint(len(r.alerts) + 1)
	r.alerts = append(r.alerts, *alert)
	return nil
}

func (r *AlertRepo) FindByUserID(userID uint) ([]model.RateAlert, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.RateAlert
	for _, a := range r.alerts {
		if a.UserID == userID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *AlertRepo) FindUntriggered() ([]model.RateAlert, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.RateAlert
	for _, a := range r.alerts {
		if !a.Triggered {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *AlertRepo) FindUntriggeredByPair(from, to string) ([]model.RateAlert, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.RateAlert
	for _, a := range r.alerts {
		if !a.Triggered && a.FromCurrency == from && a.ToCurrency == to {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *AlertRepo) MarkTriggered(id uint) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, a := range r.alerts {
		if a.ID == id {
			r.alerts[i].Triggered = true
			return nil
		}
	}
	return nil
}

func (r *AlertRepo) Delete(id uint, userID uint) (bool, error) {
	if r.Err != nil {
		return false, r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []model.RateAlert
	found := false
	for _, a := range r.alerts {
		if a.ID == id && a.UserID == userID {
			found = true
		} else {
			result = append(result, a)
		}
	}
	r.alerts = result
	return found, nil
}

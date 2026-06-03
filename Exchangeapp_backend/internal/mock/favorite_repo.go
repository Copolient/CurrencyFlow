package mock

import (
	"exchangeapp/internal/model"
	"sync"
)

type FavoriteRepo struct {
	mu        sync.RWMutex
	favorites []model.Favorite
	Err       error
}

func NewFavoriteRepo() *FavoriteRepo {
	return &FavoriteRepo{}
}

func (r *FavoriteRepo) Create(favorite *model.Favorite) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.favorites = append(r.favorites, *favorite)
	return nil
}

func (r *FavoriteRepo) FindByUserID(userID uint) ([]model.Favorite, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Favorite
	for _, f := range r.favorites {
		if f.UserID == userID {
			result = append(result, f)
		}
	}
	return result, nil
}

func (r *FavoriteRepo) Delete(userID uint, from, to string) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []model.Favorite
	for _, f := range r.favorites {
		if !(f.UserID == userID && f.FromCurrency == from && f.ToCurrency == to) {
			result = append(result, f)
		}
	}
	r.favorites = result
	return nil
}

func (r *FavoriteRepo) Exists(userID uint, from, to string) (bool, error) {
	if r.Err != nil {
		return false, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, f := range r.favorites {
		if f.UserID == userID && f.FromCurrency == from && f.ToCurrency == to {
			return true, nil
		}
	}
	return false, nil
}

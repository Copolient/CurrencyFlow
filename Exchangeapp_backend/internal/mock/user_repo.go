package mock

import (
	"exchangeapp/internal/model"
	"sync"
)

type UserRepo struct {
	mu    sync.RWMutex
	users map[string]*model.User
	Err   error
}

func NewUserRepo() *UserRepo {
	return &UserRepo{users: make(map[string]*model.User)}
}

func (r *UserRepo) Create(user *model.User) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.Username] = user
	return nil
}

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[username]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}

func (r *UserRepo) FindByID(id uint) (*model.User, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

func (r *UserRepo) FindByIDs(ids []uint) ([]model.User, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	idSet := make(map[uint]bool)
	for _, id := range ids {
		idSet[id] = true
	}
	var result []model.User
	for _, u := range r.users {
		if idSet[u.ID] {
			result = append(result, *u)
		}
	}
	return result, nil
}

func (r *UserRepo) Update(user *model.User) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.Username] = user
	return nil
}

func (r *UserRepo) IncrementFollowingCount(_ uint) error {
	return nil // no-op in mock
}

func (r *UserRepo) DecrementFollowingCount(_ uint) error {
	return nil
}

func (r *UserRepo) IncrementFollowersCount(_ uint) error {
	return nil
}

func (r *UserRepo) DecrementFollowersCount(_ uint) error {
	return nil
}

var ErrNotFound = &notFoundError{}

type notFoundError struct{}

func (e *notFoundError) Error() string { return "record not found" }

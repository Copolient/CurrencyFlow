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

var ErrNotFound = &notFoundError{}

type notFoundError struct{}

func (e *notFoundError) Error() string { return "record not found" }

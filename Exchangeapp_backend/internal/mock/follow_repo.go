package mock

import (
	"exchangeapp/internal/model"
	"sync"
)

type FollowRepo struct {
	mu      sync.RWMutex
	follows []model.Follow
	Err     error
}

func NewFollowRepo() *FollowRepo {
	return &FollowRepo{}
}

func (r *FollowRepo) Create(follow *model.Follow) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.follows = append(r.follows, *follow)
	return nil
}

func (r *FollowRepo) Delete(followerID, followeeID uint) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []model.Follow
	for _, f := range r.follows {
		if !(f.FollowerID == followerID && f.FolloweeID == followeeID) {
			result = append(result, f)
		}
	}
	r.follows = result
	return nil
}

func (r *FollowRepo) Exists(followerID, followeeID uint) (bool, error) {
	if r.Err != nil {
		return false, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, f := range r.follows {
		if f.FollowerID == followerID && f.FolloweeID == followeeID {
			return true, nil
		}
	}
	return false, nil
}

func (r *FollowRepo) FindFollowing(userID uint) ([]uint, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var ids []uint
	for _, f := range r.follows {
		if f.FollowerID == userID {
			ids = append(ids, f.FolloweeID)
		}
	}
	return ids, nil
}

func (r *FollowRepo) FindFollowers(userID uint) ([]uint, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var ids []uint
	for _, f := range r.follows {
		if f.FolloweeID == userID {
			ids = append(ids, f.FollowerID)
		}
	}
	return ids, nil
}

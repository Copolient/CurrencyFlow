package mock

import (
	"exchangeapp/internal/model"
	"sync"
)

type PostRepo struct {
	mu    sync.RWMutex
	posts []model.Post
	Err   error
}

func NewPostRepo() *PostRepo {
	return &PostRepo{}
}

func (r *PostRepo) Create(post *model.Post) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	post.ID = uint(len(r.posts) + 1)
	r.posts = append(r.posts, *post)
	return nil
}

func (r *PostRepo) FindAll(limit, offset int) ([]model.Post, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	start := offset
	if start > len(r.posts) {
		return nil, nil
	}
	end := start + limit
	if end > len(r.posts) {
		end = len(r.posts)
	}
	return r.posts[start:end], nil
}

func (r *PostRepo) FindByUserID(userID uint, limit, offset int) ([]model.Post, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Post
	for _, p := range r.posts {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	start := offset
	if start > len(result) {
		return nil, nil
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], nil
}

func (r *PostRepo) FindByFollowing(userIDs []uint, limit, offset int) ([]model.Post, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	idSet := make(map[uint]bool)
	for _, id := range userIDs {
		idSet[id] = true
	}
	var result []model.Post
	for _, p := range r.posts {
		if idSet[p.UserID] {
			result = append(result, p)
		}
	}
	start := offset
	if start > len(result) {
		return nil, nil
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], nil
}

func (r *PostRepo) FindByID(id uint) (*model.Post, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i, p := range r.posts {
		if p.ID == id {
			return &r.posts[i], nil
		}
	}
	return nil, ErrNotFound
}

func (r *PostRepo) IncrementLikes(id uint) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, p := range r.posts {
		if p.ID == id {
			r.posts[i].Likes++
			return nil
		}
	}
	return nil
}

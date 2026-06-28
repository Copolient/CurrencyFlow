package mock

import (
	"exchangeapp/internal/model"
	"sync"
	"sync/atomic"
)

type ArticleRepo struct {
	mu       sync.RWMutex
	articles map[string]*model.Article
	counter  atomic.Int64
	Err      error
}

func NewArticleRepo() *ArticleRepo {
	return &ArticleRepo{articles: make(map[string]*model.Article)}
}

func (r *ArticleRepo) Create(article *model.Article) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter.Add(1)
	r.articles[article.Title] = article
	return nil
}

func (r *ArticleRepo) FindAll() ([]model.Article, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]model.Article, 0, len(r.articles))
	for _, a := range r.articles {
		result = append(result, *a)
	}
	return result, nil
}

func (r *ArticleRepo) FindAllPaginated(limit, offset int) ([]model.Article, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]model.Article, 0, len(r.articles))
	for _, a := range r.articles {
		result = append(result, *a)
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

func (r *ArticleRepo) FindByID(id string) (*model.Article, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, a := range r.articles {
		if a.Title == id { // simplified: use title as key for test
			return a, nil
		}
	}
	return nil, ErrNotFound
}

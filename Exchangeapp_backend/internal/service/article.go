package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"exchangeapp/internal/model"
	"exchangeapp/internal/repository"
	"exchangeapp/pkg/cache"
	"time"

	"gorm.io/gorm"
)

const articleCacheKey = "articles"

type ArticleService struct {
	articleRepo repository.ArticleRepository
	cache       cache.Cache
}

func NewArticleService(articleRepo repository.ArticleRepository, cache cache.Cache) *ArticleService {
	return &ArticleService{articleRepo: articleRepo, cache: cache}
}

func (s *ArticleService) CreateArticle(ctx context.Context, article *model.Article) error {
	if err := s.articleRepo.Create(article); err != nil {
		return fmt.Errorf("articleRepo.Create: %w", err)
	}
	_ = s.cache.Del(ctx, articleCacheKey)
	return nil
}

func (s *ArticleService) GetArticles(ctx context.Context) ([]model.Article, error) {
	cached, err := s.cache.Get(ctx, articleCacheKey)
	if err == nil {
		var articles []model.Article
		if json.Unmarshal([]byte(cached), &articles) == nil {
			return articles, nil
		}
	}

	articles, err := s.articleRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("articleRepo.FindAll: %w", err)
	}

	if data, err := json.Marshal(articles); err == nil {
		_ = s.cache.Set(ctx, articleCacheKey, string(data), 10*time.Minute)
	}

	return articles, nil
}

func (s *ArticleService) GetArticleByID(ctx context.Context, id string) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("articleRepo.FindByID(%s): %w", id, err)
	}
	return article, nil
}

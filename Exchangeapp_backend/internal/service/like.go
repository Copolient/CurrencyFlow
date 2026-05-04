package service

import (
	"context"
	"fmt"
	"exchangeapp/pkg/cache"
)

type LikeService struct {
	cache cache.Cache
}

func NewLikeService(cache cache.Cache) *LikeService {
	return &LikeService{cache: cache}
}

func (s *LikeService) LikeArticle(ctx context.Context, articleID string) error {
	key := fmt.Sprintf("article:%s:likes", articleID)
	if _, err := s.cache.Incr(ctx, key); err != nil {
		return fmt.Errorf("cache.Incr(%s): %w", key, err)
	}
	return nil
}

func (s *LikeService) GetArticleLikes(ctx context.Context, articleID string) (string, error) {
	key := fmt.Sprintf("article:%s:likes", articleID)
	likes, err := s.cache.Get(ctx, key)
	if err != nil {
		return "0", nil
	}
	return likes, nil
}

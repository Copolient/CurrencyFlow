package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"exchangeapp/pkg/cache"
)

type LikeService struct {
	cache cache.Cache
}

func NewLikeService(cache cache.Cache) *LikeService {
	return &LikeService{cache: cache}
}

func (s *LikeService) LikeArticle(ctx context.Context, articleID string, userID uint) (bool, error) {
	// Check if user already liked this article
	likedKey := fmt.Sprintf("article:%s:liked:%d", articleID, userID)
	exists, err := s.cache.Get(ctx, likedKey)
	if err == nil && exists == "1" {
		return false, nil // already liked
	}

	// Record user's like
	if err := s.cache.Set(ctx, likedKey, "1", 30*24*time.Hour); err != nil {
		return false, fmt.Errorf("cache.Set(%s): %w", likedKey, err)
	}

	// Increment like count
	countKey := fmt.Sprintf("article:%s:likes", articleID)
	if _, err := s.cache.Incr(ctx, countKey); err != nil {
		return false, fmt.Errorf("cache.Incr(%s): %w", countKey, err)
	}

	// Set TTL on the count key if it's new
	s.cache.Expire(ctx, countKey, 30*24*time.Hour)

	return true, nil
}

func (s *LikeService) GetArticleLikes(ctx context.Context, articleID string) (int64, error) {
	key := fmt.Sprintf("article:%s:likes", articleID)
	val, err := s.cache.Get(ctx, key)
	if err != nil {
		return 0, nil
	}
	count, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, nil
	}
	return count, nil
}

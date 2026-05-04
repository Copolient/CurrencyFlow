package service_test

import (
	"context"
	"exchangeapp/internal/mock"
	"exchangeapp/internal/service"
	"testing"
)

func TestLikeArticle_Success(t *testing.T) {
	cache := mock.NewCache()
	svc := service.NewLikeService(cache)

	err := svc.LikeArticle(context.Background(), "42")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetArticleLikes_ReturnsCount(t *testing.T) {
	cache := mock.NewCache()
	svc := service.NewLikeService(cache)

	// First like
	_ = svc.LikeArticle(context.Background(), "42")

	likes, err := svc.GetArticleLikes(context.Background(), "42")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if likes == "" {
		t.Fatal("expected non-empty likes")
	}
}

func TestGetArticleLikes_NoLikes(t *testing.T) {
	cache := mock.NewCache()
	svc := service.NewLikeService(cache)

	likes, err := svc.GetArticleLikes(context.Background(), "999")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if likes != "0" {
		t.Fatalf("expected '0', got '%s'", likes)
	}
}

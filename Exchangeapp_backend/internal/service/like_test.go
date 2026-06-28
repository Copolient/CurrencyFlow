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

	liked, err := svc.LikeArticle(context.Background(), "42", 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !liked {
		t.Fatal("expected liked to be true")
	}
}

func TestLikeArticle_Duplicate(t *testing.T) {
	cache := mock.NewCache()
	svc := service.NewLikeService(cache)

	liked1, err := svc.LikeArticle(context.Background(), "42", 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !liked1 {
		t.Fatal("expected first like to be true")
	}

	liked2, err := svc.LikeArticle(context.Background(), "42", 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if liked2 {
		t.Fatal("expected second like to be false (duplicate)")
	}
}

func TestGetArticleLikes_ReturnsCount(t *testing.T) {
	cache := mock.NewCache()
	svc := service.NewLikeService(cache)

	_, _ = svc.LikeArticle(context.Background(), "42", 1)

	likes, err := svc.GetArticleLikes(context.Background(), "42")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if likes != 1 {
		t.Fatalf("expected 1, got %d", likes)
	}
}

func TestGetArticleLikes_NoLikes(t *testing.T) {
	cache := mock.NewCache()
	svc := service.NewLikeService(cache)

	likes, err := svc.GetArticleLikes(context.Background(), "999")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if likes != 0 {
		t.Fatalf("expected 0, got %d", likes)
	}
}

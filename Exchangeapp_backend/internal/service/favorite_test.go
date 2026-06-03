package service_test

import (
	"errors"
	"testing"

	"exchangeapp/internal/mock"
	"exchangeapp/internal/service"
)

func TestAddFavorite_Success(t *testing.T) {
	repo := mock.NewFavoriteRepo()
	svc := service.NewFavoriteService(repo)

	err := svc.AddFavorite(1, "USD", "CNY")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAddFavorite_AlreadyExists(t *testing.T) {
	repo := mock.NewFavoriteRepo()
	svc := service.NewFavoriteService(repo)

	// Add once
	_ = svc.AddFavorite(1, "USD", "CNY")

	// Add again - should not error
	err := svc.AddFavorite(1, "USD", "CNY")
	if err != nil {
		t.Fatalf("expected no error for duplicate, got %v", err)
	}
}

func TestAddFavorite_RepoError(t *testing.T) {
	repo := mock.NewFavoriteRepo()
	repo.Err = errors.New("db error")
	svc := service.NewFavoriteService(repo)

	err := svc.AddFavorite(1, "USD", "CNY")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetFavorites_Success(t *testing.T) {
	repo := mock.NewFavoriteRepo()
	svc := service.NewFavoriteService(repo)

	_ = svc.AddFavorite(1, "USD", "CNY")
	_ = svc.AddFavorite(1, "EUR", "USD")
	_ = svc.AddFavorite(2, "GBP", "JPY") // different user

	favorites, err := svc.GetFavorites(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(favorites) != 2 {
		t.Fatalf("expected 2 favorites, got %d", len(favorites))
	}
}

func TestGetFavorites_Empty(t *testing.T) {
	repo := mock.NewFavoriteRepo()
	svc := service.NewFavoriteService(repo)

	favorites, err := svc.GetFavorites(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(favorites) != 0 {
		t.Fatalf("expected 0 favorites, got %d", len(favorites))
	}
}

func TestGetFavorites_RepoError(t *testing.T) {
	repo := mock.NewFavoriteRepo()
	repo.Err = errors.New("db error")
	svc := service.NewFavoriteService(repo)

	_, err := svc.GetFavorites(1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRemoveFavorite_Success(t *testing.T) {
	repo := mock.NewFavoriteRepo()
	svc := service.NewFavoriteService(repo)

	_ = svc.AddFavorite(1, "USD", "CNY")
	_ = svc.AddFavorite(1, "EUR", "USD")

	err := svc.RemoveFavorite(1, "USD", "CNY")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	favorites, _ := svc.GetFavorites(1)
	if len(favorites) != 1 {
		t.Fatalf("expected 1 favorite after removal, got %d", len(favorites))
	}
	if favorites[0].FromCurrency != "EUR" {
		t.Fatalf("expected remaining favorite to be EUR, got %s", favorites[0].FromCurrency)
	}
}

func TestRemoveFavorite_NotFound(t *testing.T) {
	repo := mock.NewFavoriteRepo()
	svc := service.NewFavoriteService(repo)

	// Remove non-existent - should not error
	err := svc.RemoveFavorite(1, "USD", "CNY")
	if err != nil {
		t.Fatalf("expected no error for non-existent, got %v", err)
	}
}

func TestRemoveFavorite_RepoError(t *testing.T) {
	repo := mock.NewFavoriteRepo()
	repo.Err = errors.New("db error")
	svc := service.NewFavoriteService(repo)

	err := svc.RemoveFavorite(1, "USD", "CNY")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

package service

import (
	"fmt"

	"exchangeapp/internal/model"
	"exchangeapp/internal/repository"
)

type FavoriteService struct {
	repo repository.FavoriteRepository
}

func NewFavoriteService(repo repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo}
}

func (s *FavoriteService) AddFavorite(userID uint, from, to string) error {
	exists, err := s.repo.Exists(userID, from, to)
	if err != nil {
		return fmt.Errorf("check favorite exists: %w", err)
	}
	if exists {
		return nil // already favorited
	}

	fav := &model.Favorite{
		UserID:       userID,
		FromCurrency: from,
		ToCurrency:   to,
	}
	if err := s.repo.Create(fav); err != nil {
		return fmt.Errorf("favoriteRepo.Create: %w", err)
	}
	return nil
}

func (s *FavoriteService) GetFavorites(userID uint) ([]model.Favorite, error) {
	favorites, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("favoriteRepo.FindByUserID: %w", err)
	}
	return favorites, nil
}

func (s *FavoriteService) RemoveFavorite(userID uint, from, to string) (bool, error) {
	exists, err := s.repo.Exists(userID, from, to)
	if err != nil {
		return false, fmt.Errorf("check favorite exists: %w", err)
	}
	if !exists {
		return false, nil
	}

	if err := s.repo.Delete(userID, from, to); err != nil {
		return false, fmt.Errorf("favoriteRepo.Delete: %w", err)
	}
	return true, nil
}

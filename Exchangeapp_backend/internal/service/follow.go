package service

import (
	"fmt"

	"exchangeapp/internal/model"
	"exchangeapp/internal/repository"
)

type FollowService struct {
	followRepo repository.FollowRepository
	userRepo   repository.UserRepository
}

func NewFollowService(followRepo repository.FollowRepository, userRepo repository.UserRepository) *FollowService {
	return &FollowService{
		followRepo: followRepo,
		userRepo:   userRepo,
	}
}

func (s *FollowService) Follow(followerID, followeeID uint) error {
	if followerID == followeeID {
		return fmt.Errorf("cannot follow yourself")
	}

	exists, err := s.followRepo.Exists(followerID, followeeID)
	if err != nil {
		return fmt.Errorf("check follow exists: %w", err)
	}
	if exists {
		return nil // already following
	}

	follow := &model.Follow{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	if err := s.followRepo.Create(follow); err != nil {
		return fmt.Errorf("followRepo.Create: %w", err)
	}

	// Use atomic count updates to avoid race conditions
	_ = s.userRepo.IncrementFollowingCount(followerID)
	_ = s.userRepo.IncrementFollowersCount(followeeID)

	return nil
}

func (s *FollowService) Unfollow(followerID, followeeID uint) error {
	deleted, err := s.followRepo.Delete(followerID, followeeID)
	if err != nil {
		return fmt.Errorf("followRepo.Delete: %w", err)
	}

	if deleted {
		// Use atomic count updates
		_ = s.userRepo.DecrementFollowingCount(followerID)
		_ = s.userRepo.DecrementFollowersCount(followeeID)
	}

	return nil
}

func (s *FollowService) IsFollowing(followerID, followeeID uint) (bool, error) {
	exists, err := s.followRepo.Exists(followerID, followeeID)
	if err != nil {
		return false, fmt.Errorf("followRepo.Exists: %w", err)
	}
	return exists, nil
}

func (s *FollowService) GetFollowing(userID uint) ([]uint, error) {
	ids, err := s.followRepo.FindFollowing(userID)
	if err != nil {
		return nil, fmt.Errorf("followRepo.FindFollowing: %w", err)
	}
	return ids, nil
}

func (s *FollowService) GetFollowers(userID uint) ([]uint, error) {
	ids, err := s.followRepo.FindFollowers(userID)
	if err != nil {
		return nil, fmt.Errorf("followRepo.FindFollowers: %w", err)
	}
	return ids, nil
}

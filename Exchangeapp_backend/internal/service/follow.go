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

	// Update counts
	s.updateFollowCounts(followerID, followeeID)

	return nil
}

func (s *FollowService) Unfollow(followerID, followeeID uint) error {
	if err := s.followRepo.Delete(followerID, followeeID); err != nil {
		return fmt.Errorf("followRepo.Delete: %w", err)
	}

	// Update counts
	s.updateFollowCounts(followerID, followeeID)

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

func (s *FollowService) updateFollowCounts(followerID, followeeID uint) {
	// Update follower's following count
	following, _ := s.followRepo.FindFollowing(followerID)
	if user, err := s.userRepo.FindByID(followerID); err == nil {
		user.FollowingCount = len(following)
		_ = s.userRepo.Update(user)
	}

	// Update followee's followers count
	followers, _ := s.followRepo.FindFollowers(followeeID)
	if user, err := s.userRepo.FindByID(followeeID); err == nil {
		user.FollowersCount = len(followers)
		_ = s.userRepo.Update(user)
	}
}

package service_test

import (
	"testing"

	"exchangeapp/internal/mock"
	"exchangeapp/internal/model"
	"exchangeapp/internal/service"
)

func TestFollow_Success(t *testing.T) {
	followRepo := mock.NewFollowRepo()
	userRepo := mock.NewUserRepo()
	svc := service.NewFollowService(followRepo, userRepo)

	// Create users
	userRepo.Create(&model.User{Username: "alice"})
	userRepo.Create(&model.User{Username: "bob"})

	err := svc.Follow(1, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFollow_Self(t *testing.T) {
	followRepo := mock.NewFollowRepo()
	userRepo := mock.NewUserRepo()
	svc := service.NewFollowService(followRepo, userRepo)

	err := svc.Follow(1, 1)
	if err == nil {
		t.Fatal("expected error for self-follow, got nil")
	}
}

func TestFollow_AlreadyFollowing(t *testing.T) {
	followRepo := mock.NewFollowRepo()
	userRepo := mock.NewUserRepo()
	svc := service.NewFollowService(followRepo, userRepo)

	_ = svc.Follow(1, 2)

	// Follow again - should not error
	err := svc.Follow(1, 2)
	if err != nil {
		t.Fatalf("expected no error for duplicate follow, got %v", err)
	}
}

func TestUnfollow_Success(t *testing.T) {
	followRepo := mock.NewFollowRepo()
	userRepo := mock.NewUserRepo()
	svc := service.NewFollowService(followRepo, userRepo)

	_ = svc.Follow(1, 2)

	err := svc.Unfollow(1, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	isFollowing, _ := svc.IsFollowing(1, 2)
	if isFollowing {
		t.Fatal("expected not following after unfollow")
	}
}

func TestIsFollowing_True(t *testing.T) {
	followRepo := mock.NewFollowRepo()
	userRepo := mock.NewUserRepo()
	svc := service.NewFollowService(followRepo, userRepo)

	_ = svc.Follow(1, 2)

	isFollowing, err := svc.IsFollowing(1, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !isFollowing {
		t.Fatal("expected to be following")
	}
}

func TestIsFollowing_False(t *testing.T) {
	followRepo := mock.NewFollowRepo()
	userRepo := mock.NewUserRepo()
	svc := service.NewFollowService(followRepo, userRepo)

	isFollowing, err := svc.IsFollowing(1, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if isFollowing {
		t.Fatal("expected not to be following")
	}
}

func TestGetFollowing_Success(t *testing.T) {
	followRepo := mock.NewFollowRepo()
	userRepo := mock.NewUserRepo()
	svc := service.NewFollowService(followRepo, userRepo)

	_ = svc.Follow(1, 2)
	_ = svc.Follow(1, 3)

	following, err := svc.GetFollowing(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(following) != 2 {
		t.Fatalf("expected 2 following, got %d", len(following))
	}
}

func TestGetFollowers_Success(t *testing.T) {
	followRepo := mock.NewFollowRepo()
	userRepo := mock.NewUserRepo()
	svc := service.NewFollowService(followRepo, userRepo)

	_ = svc.Follow(1, 3)
	_ = svc.Follow(2, 3)

	followers, err := svc.GetFollowers(3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(followers) != 2 {
		t.Fatalf("expected 2 followers, got %d", len(followers))
	}
}

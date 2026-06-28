package service_test

import (
	"testing"

	"exchangeapp/internal/mock"
	"exchangeapp/internal/model"
	"exchangeapp/internal/service"
)

func TestCreatePost_Success(t *testing.T) {
	postRepo := mock.NewPostRepo()
	userRepo := mock.NewUserRepo()
	followRepo := mock.NewFollowRepo()
	svc := service.NewPostService(postRepo, userRepo, followRepo)

	err := svc.CreatePost(1, "USD/CNY looks bullish!", "USD/CNY")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetPosts_Latest(t *testing.T) {
	postRepo := mock.NewPostRepo()
	userRepo := mock.NewUserRepo()
	followRepo := mock.NewFollowRepo()
	svc := service.NewPostService(postRepo, userRepo, followRepo)

	_ = svc.CreatePost(1, "Post 1", "")
	_ = svc.CreatePost(2, "Post 2", "")
	_ = svc.CreatePost(3, "Post 3", "")

	posts, err := svc.GetPosts("latest", 0, 1, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(posts))
	}
}

func TestGetPosts_Following(t *testing.T) {
	postRepo := mock.NewPostRepo()
	userRepo := mock.NewUserRepo()
	followRepo := mock.NewFollowRepo()
	svc := service.NewPostService(postRepo, userRepo, followRepo)

	// User 1 follows user 2
	followRepo.Create(&model.Follow{FollowerID: 1, FolloweeID: 2})

	_ = svc.CreatePost(1, "My post", "")
	_ = svc.CreatePost(2, "Followed post", "")
	_ = svc.CreatePost(3, "Other post", "")

	posts, err := svc.GetPosts("following", 1, 1, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post from following, got %d", len(posts))
	}
	if posts[0].Content != "Followed post" {
		t.Fatalf("expected 'Followed post', got '%s'", posts[0].Content)
	}
}

func TestGetPosts_User(t *testing.T) {
	postRepo := mock.NewPostRepo()
	userRepo := mock.NewUserRepo()
	followRepo := mock.NewFollowRepo()
	svc := service.NewPostService(postRepo, userRepo, followRepo)

	_ = svc.CreatePost(1, "My post 1", "")
	_ = svc.CreatePost(1, "My post 2", "")
	_ = svc.CreatePost(2, "Other post", "")

	posts, err := svc.GetPosts("user", 1, 1, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
}

func TestLikePost_Success(t *testing.T) {
	postRepo := mock.NewPostRepo()
	userRepo := mock.NewUserRepo()
	followRepo := mock.NewFollowRepo()
	svc := service.NewPostService(postRepo, userRepo, followRepo)

	_ = svc.CreatePost(1, "Test post", "")

	liked, err := svc.LikePost(1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !liked {
		t.Fatal("expected liked to be true")
	}

	post, _ := svc.FindByID(1)
	if post.Likes != 1 {
		t.Fatalf("expected 1 like, got %d", post.Likes)
	}
}

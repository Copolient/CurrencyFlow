package service

import (
	"fmt"

	"exchangeapp/internal/model"
	"exchangeapp/internal/repository"
)

type PostService struct {
	postRepo   repository.PostRepository
	userRepo   repository.UserRepository
	followRepo repository.FollowRepository
}

func NewPostService(
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	followRepo repository.FollowRepository,
) *PostService {
	return &PostService{
		postRepo:   postRepo,
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

func (s *PostService) CreatePost(userID uint, content, currency string) error {
	post := &model.Post{
		UserID:   userID,
		Content:  content,
		Currency: currency,
	}
	if err := s.postRepo.Create(post); err != nil {
		return fmt.Errorf("postRepo.Create: %w", err)
	}
	return nil
}

func (s *PostService) GetPosts(feedType string, userID uint, page, pageSize int) ([]model.Post, error) {
	offset := (page - 1) * pageSize

	var posts []model.Post
	var err error

	switch feedType {
	case "following":
		followingIDs, err := s.followRepo.FindFollowing(userID)
		if err != nil {
			return nil, fmt.Errorf("followRepo.FindFollowing: %w", err)
		}
		posts, err = s.postRepo.FindByFollowing(followingIDs, pageSize, offset)
	case "user":
		posts, err = s.postRepo.FindByUserID(userID, pageSize, offset)
	default: // "latest" or "hot"
		posts, err = s.postRepo.FindAll(pageSize, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("postRepo.Find: %w", err)
	}

	// Enrich with username
	for i := range posts {
		user, err := s.userRepo.FindByID(posts[i].UserID)
		if err == nil && user != nil {
			posts[i].Username = user.Username
		}
	}

	return posts, nil
}

func (s *PostService) LikePost(id uint) error {
	if err := s.postRepo.IncrementLikes(id); err != nil {
		return fmt.Errorf("postRepo.IncrementLikes: %w", err)
	}
	return nil
}

// FindByID is a helper needed for enriching posts
func (s *PostService) FindByID(id uint) (*model.Post, error) {
	return s.postRepo.FindByID(id)
}

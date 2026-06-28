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

	if len(posts) == 0 {
		return posts, nil
	}

	// Batch fetch users to avoid N+1 query
	userIDs := make([]uint, 0, len(posts))
	seen := make(map[uint]bool)
	for _, p := range posts {
		if !seen[p.UserID] {
			seen[p.UserID] = true
			userIDs = append(userIDs, p.UserID)
		}
	}

	users, err := s.userRepo.FindByIDs(userIDs)
	if err == nil {
		userMap := make(map[uint]string, len(users))
		for _, u := range users {
			userMap[u.ID] = u.Username
		}
		for i := range posts {
			if name, ok := userMap[posts[i].UserID]; ok {
				posts[i].Username = name
			}
		}
	}

	return posts, nil
}

func (s *PostService) LikePost(postID, userID uint) (bool, error) {
	// Check if already liked
	liked, err := s.postRepo.HasUserLiked(postID, userID)
	if err != nil {
		return false, fmt.Errorf("postRepo.HasUserLiked: %w", err)
	}
	if liked {
		return false, nil // already liked
	}

	// Record the like and increment count atomically
	if err := s.postRepo.AddLike(postID, userID); err != nil {
		return false, fmt.Errorf("postRepo.AddLike: %w", err)
	}

	return true, nil
}

// FindByID is a helper needed for enriching posts
func (s *PostService) FindByID(id uint) (*model.Post, error) {
	return s.postRepo.FindByID(id)
}

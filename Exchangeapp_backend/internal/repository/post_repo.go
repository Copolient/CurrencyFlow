package repository

import (
	"exchangeapp/internal/model"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *model.Post) error
	FindAll(limit, offset int) ([]model.Post, error)
	FindByUserID(userID uint, limit, offset int) ([]model.Post, error)
	FindByFollowing(userIDs []uint, limit, offset int) ([]model.Post, error)
	FindByID(id uint) (*model.Post, error)
	IncrementLikes(id uint) error
	HasUserLiked(postID, userID uint) (bool, error)
	AddLike(postID, userID uint) error
}

type postRepo struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepo{db: db}
}

func (r *postRepo) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepo) FindAll(limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts).Error
	return posts, err
}

func (r *postRepo) FindByUserID(userID uint, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts).Error
	return posts, err
}

func (r *postRepo) FindByFollowing(userIDs []uint, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	if len(userIDs) == 0 {
		return posts, nil
	}
	err := r.db.Where("user_id IN ?", userIDs).Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts).Error
	return posts, err
}

func (r *postRepo) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepo) IncrementLikes(id uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).Update("likes", gorm.Expr("likes + 1")).Error
}

func (r *postRepo) HasUserLiked(postID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.PostLike{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *postRepo) AddLike(postID, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Insert the like record
		like := &model.PostLike{
			PostID: postID,
			UserID: userID,
		}
		if err := tx.Create(like).Error; err != nil {
			return err
		}
		// Increment the post's like count
		return tx.Model(&model.Post{}).Where("id = ?", postID).
			Update("likes", gorm.Expr("likes + 1")).Error
	})
}

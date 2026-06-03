package repository

import (
	"exchangeapp/internal/model"

	"gorm.io/gorm"
)

type FollowRepository interface {
	Create(follow *model.Follow) error
	Delete(followerID, followeeID uint) error
	Exists(followerID, followeeID uint) (bool, error)
	FindFollowing(userID uint) ([]uint, error)
	FindFollowers(userID uint) ([]uint, error)
}

type followRepo struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepo{db: db}
}

func (r *followRepo) Create(follow *model.Follow) error {
	return r.db.Create(follow).Error
}

func (r *followRepo) Delete(followerID, followeeID uint) error {
	return r.db.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&model.Follow{}).Error
}

func (r *followRepo) Exists(followerID, followeeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).Error
	return count > 0, err
}

func (r *followRepo) FindFollowing(userID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&model.Follow{}).Where("follower_id = ?", userID).Pluck("followee_id", &ids).Error
	return ids, err
}

func (r *followRepo) FindFollowers(userID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&model.Follow{}).Where("followee_id = ?", userID).Pluck("follower_id", &ids).Error
	return ids, err
}

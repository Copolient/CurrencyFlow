package repository

import (
	"exchangeapp/internal/model"

	"gorm.io/gorm"
)

type FavoriteRepository interface {
	Create(favorite *model.Favorite) error
	FindByUserID(userID uint) ([]model.Favorite, error)
	Delete(userID uint, from, to string) error
	Exists(userID uint, from, to string) (bool, error)
}

type favoriteRepo struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepo{db: db}
}

func (r *favoriteRepo) Create(favorite *model.Favorite) error {
	return r.db.Create(favorite).Error
}

func (r *favoriteRepo) FindByUserID(userID uint) ([]model.Favorite, error) {
	var favorites []model.Favorite
	err := r.db.Where("user_id = ?", userID).Find(&favorites).Error
	return favorites, err
}

func (r *favoriteRepo) Delete(userID uint, from, to string) error {
	return r.db.Where("user_id = ? AND from_currency = ? AND to_currency = ?", userID, from, to).
		Delete(&model.Favorite{}).Error
}

func (r *favoriteRepo) Exists(userID uint, from, to string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Favorite{}).
		Where("user_id = ? AND from_currency = ? AND to_currency = ?", userID, from, to).
		Count(&count).Error
	return count > 0, err
}

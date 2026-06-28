package repository

import (
	"exchangeapp/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindByIDs(ids []uint) ([]model.User, error)
	Update(user *model.User) error
	IncrementFollowingCount(userID uint) error
	DecrementFollowingCount(userID uint) error
	IncrementFollowersCount(userID uint) error
	DecrementFollowersCount(userID uint) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByIDs(ids []uint) ([]model.User, error) {
	var users []model.User
	if len(ids) == 0 {
		return users, nil
	}
	err := r.db.Where("id IN ?", ids).Find(&users).Error
	return users, err
}

func (r *userRepo) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) IncrementFollowingCount(userID uint) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Update("following_count", gorm.Expr("following_count + 1")).Error
}

func (r *userRepo) DecrementFollowingCount(userID uint) error {
	return r.db.Model(&model.User{}).Where("id = ? AND following_count > 0", userID).
		Update("following_count", gorm.Expr("following_count - 1")).Error
}

func (r *userRepo) IncrementFollowersCount(userID uint) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Update("followers_count", gorm.Expr("followers_count + 1")).Error
}

func (r *userRepo) DecrementFollowersCount(userID uint) error {
	return r.db.Model(&model.User{}).Where("id = ? AND followers_count > 0", userID).
		Update("followers_count", gorm.Expr("followers_count - 1")).Error
}

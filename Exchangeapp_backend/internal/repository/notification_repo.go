package repository

import (
	"exchangeapp/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *model.Notification) error
	FindByUserID(userID uint) ([]model.Notification, error)
	MarkRead(id uint, userID uint) error
	MarkAllRead(userID uint) error
	CountUnread(userID uint) (int64, error)
}

type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepo) FindByUserID(userID uint) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(50).Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepo) MarkRead(id uint, userID uint) error {
	return r.db.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("read", true).Error
}

func (r *notificationRepo) MarkAllRead(userID uint) error {
	return r.db.Model(&model.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Update("read", true).Error
}

func (r *notificationRepo) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Count(&count).Error
	return count, err
}

package repository

import (
	"context"

	"exchangeapp/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *model.Notification) error
	FindByUserID(ctx context.Context, userID uint) ([]model.Notification, error)
	MarkRead(ctx context.Context, id uint, userID uint) error
	MarkAllRead(ctx context.Context, userID uint) error
	CountUnread(ctx context.Context, userID uint) (int64, error)
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

func (r *notificationRepo) FindByUserID(ctx context.Context, userID uint) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(50).Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepo) MarkRead(ctx context.Context, id uint, userID uint) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("read", true).Error
}

func (r *notificationRepo) MarkAllRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Update("read", true).Error
}

func (r *notificationRepo) CountUnread(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Count(&count).Error
	return count, err
}

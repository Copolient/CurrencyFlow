package service

import (
	"fmt"

	"exchangeapp/internal/model"
	"exchangeapp/internal/repository"
)

type NotificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) GetNotifications(userID uint) ([]model.Notification, error) {
	notifications, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("notificationRepo.FindByUserID: %w", err)
	}
	return notifications, nil
}

func (s *NotificationService) MarkRead(id uint, userID uint) error {
	if err := s.repo.MarkRead(id, userID); err != nil {
		return fmt.Errorf("notificationRepo.MarkRead: %w", err)
	}
	return nil
}

func (s *NotificationService) MarkAllRead(userID uint) error {
	if err := s.repo.MarkAllRead(userID); err != nil {
		return fmt.Errorf("notificationRepo.MarkAllRead: %w", err)
	}
	return nil
}

func (s *NotificationService) CountUnread(userID uint) (int64, error) {
	count, err := s.repo.CountUnread(userID)
	if err != nil {
		return 0, fmt.Errorf("notificationRepo.CountUnread: %w", err)
	}
	return count, nil
}

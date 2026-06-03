package mock

import (
	"exchangeapp/internal/model"
	"sync"
)

type NotificationRepo struct {
	mu            sync.RWMutex
	notifications []model.Notification
	Err           error
}

func NewNotificationRepo() *NotificationRepo {
	return &NotificationRepo{}
}

func (r *NotificationRepo) Create(notification *model.Notification) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	notification.ID = uint(len(r.notifications) + 1)
	r.notifications = append(r.notifications, *notification)
	return nil
}

func (r *NotificationRepo) FindByUserID(userID uint) ([]model.Notification, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Notification
	for _, n := range r.notifications {
		if n.UserID == userID {
			result = append(result, n)
		}
	}
	return result, nil
}

func (r *NotificationRepo) MarkRead(id uint, userID uint) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, n := range r.notifications {
		if n.ID == id && n.UserID == userID {
			r.notifications[i].Read = true
			return nil
		}
	}
	return nil
}

func (r *NotificationRepo) MarkAllRead(userID uint) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, n := range r.notifications {
		if n.UserID == userID {
			r.notifications[i].Read = true
		}
	}
	return nil
}

func (r *NotificationRepo) CountUnread(userID uint) (int64, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var count int64
	for _, n := range r.notifications {
		if n.UserID == userID && !n.Read {
			count++
		}
	}
	return count, nil
}

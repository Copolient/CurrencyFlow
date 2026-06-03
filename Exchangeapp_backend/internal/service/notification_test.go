package service_test

import (
	"errors"
	"testing"

	"exchangeapp/internal/mock"
	"exchangeapp/internal/model"
	"exchangeapp/internal/service"
)

func TestGetNotifications_Success(t *testing.T) {
	repo := mock.NewNotificationRepo()
	svc := service.NewNotificationService(repo)

	// Create notifications
	repo.Create(&model.Notification{UserID: 1, Title: "Test 1"})
	repo.Create(&model.Notification{UserID: 1, Title: "Test 2"})
	repo.Create(&model.Notification{UserID: 2, Title: "Other user"})

	notifications, err := svc.GetNotifications(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifications) != 2 {
		t.Fatalf("expected 2 notifications, got %d", len(notifications))
	}
}

func TestGetNotifications_Empty(t *testing.T) {
	repo := mock.NewNotificationRepo()
	svc := service.NewNotificationService(repo)

	notifications, err := svc.GetNotifications(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifications) != 0 {
		t.Fatalf("expected 0 notifications, got %d", len(notifications))
	}
}

func TestGetNotifications_RepoError(t *testing.T) {
	repo := mock.NewNotificationRepo()
	repo.Err = errors.New("db error")
	svc := service.NewNotificationService(repo)

	_, err := svc.GetNotifications(1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMarkRead_Success(t *testing.T) {
	repo := mock.NewNotificationRepo()
	svc := service.NewNotificationService(repo)

	repo.Create(&model.Notification{UserID: 1, Title: "Test"})

	err := svc.MarkRead(1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestMarkRead_RepoError(t *testing.T) {
	repo := mock.NewNotificationRepo()
	repo.Err = errors.New("db error")
	svc := service.NewNotificationService(repo)

	err := svc.MarkRead(1, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMarkAllRead_Success(t *testing.T) {
	repo := mock.NewNotificationRepo()
	svc := service.NewNotificationService(repo)

	repo.Create(&model.Notification{UserID: 1, Title: "Test 1", Read: false})
	repo.Create(&model.Notification{UserID: 1, Title: "Test 2", Read: false})

	err := svc.MarkAllRead(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	count, _ := svc.CountUnread(1)
	if count != 0 {
		t.Fatalf("expected 0 unread after mark all, got %d", count)
	}
}

func TestMarkAllRead_RepoError(t *testing.T) {
	repo := mock.NewNotificationRepo()
	repo.Err = errors.New("db error")
	svc := service.NewNotificationService(repo)

	err := svc.MarkAllRead(1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCountUnread_Success(t *testing.T) {
	repo := mock.NewNotificationRepo()
	svc := service.NewNotificationService(repo)

	repo.Create(&model.Notification{UserID: 1, Title: "Test 1", Read: false})
	repo.Create(&model.Notification{UserID: 1, Title: "Test 2", Read: true})
	repo.Create(&model.Notification{UserID: 1, Title: "Test 3", Read: false})

	count, err := svc.CountUnread(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 unread, got %d", count)
	}
}

func TestCountUnread_RepoError(t *testing.T) {
	repo := mock.NewNotificationRepo()
	repo.Err = errors.New("db error")
	svc := service.NewNotificationService(repo)

	_, err := svc.CountUnread(1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

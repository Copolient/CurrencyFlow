package service_test

import (
	"errors"
	"testing"

	"exchangeapp/internal/mock"
	"exchangeapp/internal/service"
)

func TestCreateAlert_Success(t *testing.T) {
	alertRepo := mock.NewAlertRepo()
	notifRepo := mock.NewNotificationRepo()
	svc := service.NewAlertService(alertRepo, notifRepo, nil)

	err := svc.CreateAlert(1, "USD", "CNY", 7.3, "above")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCreateAlert_RepoError(t *testing.T) {
	alertRepo := mock.NewAlertRepo()
	alertRepo.Err = errors.New("db error")
	notifRepo := mock.NewNotificationRepo()
	svc := service.NewAlertService(alertRepo, notifRepo, nil)

	err := svc.CreateAlert(1, "USD", "CNY", 7.3, "above")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUserAlerts_Success(t *testing.T) {
	alertRepo := mock.NewAlertRepo()
	notifRepo := mock.NewNotificationRepo()
	svc := service.NewAlertService(alertRepo, notifRepo, nil)

	_ = svc.CreateAlert(1, "USD", "CNY", 7.3, "above")
	_ = svc.CreateAlert(1, "EUR", "USD", 1.1, "below")
	_ = svc.CreateAlert(2, "GBP", "JPY", 150.0, "above")

	alerts, err := svc.GetUserAlerts(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}
}

func TestGetUserAlerts_Empty(t *testing.T) {
	alertRepo := mock.NewAlertRepo()
	notifRepo := mock.NewNotificationRepo()
	svc := service.NewAlertService(alertRepo, notifRepo, nil)

	alerts, err := svc.GetUserAlerts(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestGetUserAlerts_RepoError(t *testing.T) {
	alertRepo := mock.NewAlertRepo()
	alertRepo.Err = errors.New("db error")
	notifRepo := mock.NewNotificationRepo()
	svc := service.NewAlertService(alertRepo, notifRepo, nil)

	_, err := svc.GetUserAlerts(1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteAlert_Success(t *testing.T) {
	alertRepo := mock.NewAlertRepo()
	notifRepo := mock.NewNotificationRepo()
	svc := service.NewAlertService(alertRepo, notifRepo, nil)

	_ = svc.CreateAlert(1, "USD", "CNY", 7.3, "above")

	err := svc.DeleteAlert(1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	alerts, _ := svc.GetUserAlerts(1)
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts after deletion, got %d", len(alerts))
	}
}

func TestDeleteAlert_RepoError(t *testing.T) {
	alertRepo := mock.NewAlertRepo()
	alertRepo.Err = errors.New("db error")
	notifRepo := mock.NewNotificationRepo()
	svc := service.NewAlertService(alertRepo, notifRepo, nil)

	err := svc.DeleteAlert(1, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCheckAlerts_AboveTriggered(t *testing.T) {
	alertRepo := mock.NewAlertRepo()
	notifRepo := mock.NewNotificationRepo()
	svc := service.NewAlertService(alertRepo, notifRepo, nil)

	_ = svc.CreateAlert(1, "USD", "CNY", 7.3, "above")

	// Rate exceeds target
	svc.CheckAlerts("USD", "CNY", 7.35)

	// Alert should be triggered
	alerts, _ := svc.GetUserAlerts(1)
	if len(alerts) != 1 || !alerts[0].Triggered {
		t.Fatal("expected alert to be triggered")
	}

	// Notification should be created
	notifs, _ := notifRepo.FindByUserID(1)
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifs))
	}
}

func TestCheckAlerts_BelowTriggered(t *testing.T) {
	alertRepo := mock.NewAlertRepo()
	notifRepo := mock.NewNotificationRepo()
	svc := service.NewAlertService(alertRepo, notifRepo, nil)

	_ = svc.CreateAlert(1, "USD", "CNY", 7.0, "below")

	// Rate falls below target
	svc.CheckAlerts("USD", "CNY", 6.95)

	alerts, _ := svc.GetUserAlerts(1)
	if len(alerts) != 1 || !alerts[0].Triggered {
		t.Fatal("expected alert to be triggered")
	}
}

func TestCheckAlerts_NotTriggered(t *testing.T) {
	alertRepo := mock.NewAlertRepo()
	notifRepo := mock.NewNotificationRepo()
	svc := service.NewAlertService(alertRepo, notifRepo, nil)

	_ = svc.CreateAlert(1, "USD", "CNY", 7.3, "above")

	// Rate below target
	svc.CheckAlerts("USD", "CNY", 7.2)

	alerts, _ := svc.GetUserAlerts(1)
	if len(alerts) != 1 || alerts[0].Triggered {
		t.Fatal("expected alert NOT to be triggered")
	}
}

func TestCheckAlerts_WrongPair(t *testing.T) {
	alertRepo := mock.NewAlertRepo()
	notifRepo := mock.NewNotificationRepo()
	svc := service.NewAlertService(alertRepo, notifRepo, nil)

	_ = svc.CreateAlert(1, "USD", "CNY", 7.3, "above")

	// Different pair
	svc.CheckAlerts("EUR", "USD", 1.1)

	alerts, _ := svc.GetUserAlerts(1)
	if len(alerts) != 1 || alerts[0].Triggered {
		t.Fatal("expected alert NOT to be triggered for different pair")
	}
}

package service

import (
	"fmt"
	"log"

	"exchangeapp/internal/model"
	"exchangeapp/internal/repository"
	ws "exchangeapp/internal/websocket"
)

type AlertService struct {
	alertRepo      repository.AlertRepository
	notifRepo      repository.NotificationRepository
	hub            *ws.Hub
}

func NewAlertService(
	alertRepo repository.AlertRepository,
	notifRepo repository.NotificationRepository,
	hub *ws.Hub,
) *AlertService {
	return &AlertService{
		alertRepo: alertRepo,
		notifRepo: notifRepo,
		hub:       hub,
	}
}

func (s *AlertService) CreateAlert(userID uint, from, to string, targetRate float64, direction string) error {
	alert := &model.RateAlert{
		UserID:       userID,
		FromCurrency: from,
		ToCurrency:   to,
		TargetRate:   targetRate,
		Direction:    direction,
	}
	if err := s.alertRepo.Create(alert); err != nil {
		return fmt.Errorf("alertRepo.Create: %w", err)
	}
	return nil
}

func (s *AlertService) GetUserAlerts(userID uint) ([]model.RateAlert, error) {
	alerts, err := s.alertRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("alertRepo.FindByUserID: %w", err)
	}
	return alerts, nil
}

func (s *AlertService) DeleteAlert(id uint, userID uint) error {
	if err := s.alertRepo.Delete(id, userID); err != nil {
		return fmt.Errorf("alertRepo.Delete: %w", err)
	}
	return nil
}

// CheckAlerts checks all untriggered alerts against current rates
func (s *AlertService) CheckAlerts(from, to string, currentRate float64) {
	alerts, err := s.alertRepo.FindUntriggered()
	if err != nil {
		log.Printf("AlertService: failed to fetch untriggered alerts: %v", err)
		return
	}

	for _, alert := range alerts {
		if alert.FromCurrency != from || alert.ToCurrency != to {
			continue
		}

		triggered := false
		if alert.Direction == "above" && currentRate >= alert.TargetRate {
			triggered = true
		} else if alert.Direction == "below" && currentRate <= alert.TargetRate {
			triggered = true
		}

		if triggered {
			// Mark alert as triggered
			if err := s.alertRepo.MarkTriggered(alert.ID); err != nil {
				log.Printf("AlertService: failed to mark alert %d as triggered: %v", alert.ID, err)
				continue
			}

			// Create notification
			title := fmt.Sprintf("汇率预警触发: %s/%s", from, to)
			content := fmt.Sprintf("当前汇率 %.4f 已%s目标汇率 %.4f",
				currentRate,
				map[bool]string{true: "超过", false: "低于"}[alert.Direction == "above"],
				alert.TargetRate,
			)

			notif := &model.Notification{
				UserID:  alert.UserID,
				Type:    "alert_triggered",
				Title:   title,
				Content: content,
			}
			if err := s.notifRepo.Create(notif); err != nil {
				log.Printf("AlertService: failed to create notification: %v", err)
			}

			// Broadcast via WebSocket
			if s.hub != nil {
				s.hub.BroadcastRateUpdate(ws.RateUpdate{
					FromCurrency: from,
					ToCurrency:   to,
					Rate:         currentRate,
					Timestamp:    "",
				})
			}
		}
	}
}

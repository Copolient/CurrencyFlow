package repository

import (
	"exchangeapp/internal/model"

	"gorm.io/gorm"
)

type AlertRepository interface {
	Create(alert *model.RateAlert) error
	FindByUserID(userID uint) ([]model.RateAlert, error)
	FindUntriggered() ([]model.RateAlert, error)
	FindUntriggeredByPair(from, to string) ([]model.RateAlert, error)
	MarkTriggered(id uint) error
	Delete(id uint, userID uint) (bool, error)
}

type alertRepo struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) AlertRepository {
	return &alertRepo{db: db}
}

func (r *alertRepo) Create(alert *model.RateAlert) error {
	return r.db.Create(alert).Error
}

func (r *alertRepo) FindByUserID(userID uint) ([]model.RateAlert, error) {
	var alerts []model.RateAlert
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&alerts).Error
	return alerts, err
}

func (r *alertRepo) FindUntriggered() ([]model.RateAlert, error) {
	var alerts []model.RateAlert
	err := r.db.Where("triggered = ?", false).Find(&alerts).Error
	return alerts, err
}

func (r *alertRepo) FindUntriggeredByPair(from, to string) ([]model.RateAlert, error) {
	var alerts []model.RateAlert
	err := r.db.Where("triggered = ? AND from_currency = ? AND to_currency = ?", false, from, to).Find(&alerts).Error
	return alerts, err
}

func (r *alertRepo) MarkTriggered(id uint) error {
	return r.db.Model(&model.RateAlert{}).Where("id = ?", id).Update("triggered", true).Error
}

func (r *alertRepo) Delete(id uint, userID uint) (bool, error) {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.RateAlert{})
	return result.RowsAffected > 0, result.Error
}

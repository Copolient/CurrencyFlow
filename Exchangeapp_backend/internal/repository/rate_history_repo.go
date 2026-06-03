package repository

import (
	"exchangeapp/internal/model"
	"time"

	"gorm.io/gorm"
)

type RateHistoryRepository interface {
	Create(history *model.ExchangeRateHistory) error
	BulkCreate(histories []model.ExchangeRateHistory) error
	FindByPair(from, to string) ([]model.ExchangeRateHistory, error)
	FindByPairAndTimeRange(from, to string, start, end time.Time) ([]model.ExchangeRateHistory, error)
	FindLatestByAllPairs() ([]model.ExchangeRateHistory, error)
}

type rateHistoryRepo struct {
	db *gorm.DB
}

func NewRateHistoryRepository(db *gorm.DB) RateHistoryRepository {
	return &rateHistoryRepo{db: db}
}

func (r *rateHistoryRepo) Create(history *model.ExchangeRateHistory) error {
	return r.db.Create(history).Error
}

func (r *rateHistoryRepo) BulkCreate(histories []model.ExchangeRateHistory) error {
	if len(histories) == 0 {
		return nil
	}
	return r.db.Create(&histories).Error
}

func (r *rateHistoryRepo) FindByPair(from, to string) ([]model.ExchangeRateHistory, error) {
	var histories []model.ExchangeRateHistory
	err := r.db.Where("from_currency = ? AND to_currency = ?", from, to).
		Order("timestamp ASC").
		Find(&histories).Error
	return histories, err
}

func (r *rateHistoryRepo) FindByPairAndTimeRange(from, to string, start, end time.Time) ([]model.ExchangeRateHistory, error) {
	var histories []model.ExchangeRateHistory
	err := r.db.Where("from_currency = ? AND to_currency = ? AND timestamp BETWEEN ? AND ?", from, to, start, end).
		Order("timestamp ASC").
		Find(&histories).Error
	return histories, err
}

func (r *rateHistoryRepo) FindLatestByAllPairs() ([]model.ExchangeRateHistory, error) {
	var results []model.ExchangeRateHistory
	// Subquery: get max timestamp per currency pair
	subQuery := r.db.Model(&model.ExchangeRateHistory{}).
		Select("MAX(id)").
		Group("from_currency, to_currency")

	err := r.db.Where("id IN (?)", subQuery).Find(&results).Error
	return results, err
}

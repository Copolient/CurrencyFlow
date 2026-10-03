package repository

import (
	"currencyflow/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ExchangeRateRepository interface {
	Create(rate *model.ExchangeRate) error
	Upsert(rate *model.ExchangeRate) error
	FindAll() ([]model.ExchangeRate, error)
}

type exchangeRateRepo struct {
	db *gorm.DB
}

func NewExchangeRateRepository(db *gorm.DB) ExchangeRateRepository {
	return &exchangeRateRepo{db: db}
}

func (r *exchangeRateRepo) Create(rate *model.ExchangeRate) error {
	return r.db.Create(rate).Error
}

func (r *exchangeRateRepo) Upsert(rate *model.ExchangeRate) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "from_currency"}, {Name: "to_currency"}},
		DoUpdates: clause.Assignments(map[string]any{
			"rate": rate.Rate,
			"date": rate.Date,
		}),
	}).Create(rate).Error
}

func (r *exchangeRateRepo) FindAll() ([]model.ExchangeRate, error) {
	var rates []model.ExchangeRate
	err := r.db.Find(&rates).Error
	return rates, err
}

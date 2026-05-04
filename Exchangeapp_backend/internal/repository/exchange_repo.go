package repository

import (
	"exchangeapp/internal/model"

	"gorm.io/gorm"
)

type ExchangeRateRepository interface {
	Create(rate *model.ExchangeRate) error
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

func (r *exchangeRateRepo) FindAll() ([]model.ExchangeRate, error) {
	var rates []model.ExchangeRate
	err := r.db.Find(&rates).Error
	return rates, err
}

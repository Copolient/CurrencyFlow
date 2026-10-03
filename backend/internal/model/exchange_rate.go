package model

import "time"

type ExchangeRate struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	FromCurrency string    `gorm:"size:8;uniqueIndex:idx_exchange_pair" json:"fromCurrency" binding:"required,len=3"`
	ToCurrency   string    `gorm:"size:8;uniqueIndex:idx_exchange_pair" json:"toCurrency" binding:"required,len=3"`
	Rate         float64   `json:"rate" binding:"required,gt=0"`
	Date         time.Time `gorm:"index" json:"date"`
}

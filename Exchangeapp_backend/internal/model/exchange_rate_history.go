package model

import "time"

type ExchangeRateHistory struct {
	ID           uint      `gorm:"primarykey" json:"_id"`
	FromCurrency string    `gorm:"size:8;index:idx_hist_pair" json:"fromCurrency" binding:"required,len=3"`
	ToCurrency   string    `gorm:"size:8;index:idx_hist_pair" json:"toCurrency" binding:"required,len=3"`
	Rate         float64   `json:"rate" binding:"required,gt=0"`
	Timestamp    time.Time `gorm:"index" json:"timestamp"`
}

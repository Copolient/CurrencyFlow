package model

import "gorm.io/gorm"

type RateAlert struct {
	gorm.Model
	UserID       uint    `gorm:"index"`
	FromCurrency string  `gorm:"size:8" json:"fromCurrency" binding:"required,len=3"`
	ToCurrency   string  `gorm:"size:8" json:"toCurrency" binding:"required,len=3"`
	TargetRate   float64 `json:"targetRate" binding:"required,gt=0"`
	Direction    string  `gorm:"size:5" json:"direction" binding:"required,oneof=above below"`
	Triggered    bool    `gorm:"default:false" json:"triggered"`
}

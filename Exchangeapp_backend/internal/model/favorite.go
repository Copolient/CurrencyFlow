package model

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID       uint   `gorm:"uniqueIndex:idx_user_pair"`
	FromCurrency string `gorm:"size:8;uniqueIndex:idx_user_pair" json:"fromCurrency" binding:"required,len=3"`
	ToCurrency   string `gorm:"size:8;uniqueIndex:idx_user_pair" json:"toCurrency" binding:"required,len=3"`
}

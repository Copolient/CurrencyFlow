package model

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID  uint   `gorm:"index"`
	Type    string `gorm:"size:32" json:"type"`
	Title   string `gorm:"size:256" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	Read    bool   `gorm:"default:false" json:"read"`
}

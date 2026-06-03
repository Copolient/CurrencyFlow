package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID   uint   `gorm:"index" json:"userId"`
	Username string `gorm:"-" json:"username,omitempty"`
	Content  string `gorm:"type:text" json:"content" binding:"required,max=2000"`
	Currency string `gorm:"size:8" json:"currency"`
	Likes    int    `gorm:"default:0" json:"likes"`
}

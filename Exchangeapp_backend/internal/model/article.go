package model

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title   string `gorm:"size:256;index" json:"title" binding:"required,min=1,max=256"`
	Content string `gorm:"type:text" json:"content" binding:"required"`
	Preview string `gorm:"size:512" json:"preview" binding:"required,max=512"`
}

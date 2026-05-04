package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:32" json:"username" binding:"required,min=2,max=32"`
	Password string `json:"password" binding:"required,min=6"`
}

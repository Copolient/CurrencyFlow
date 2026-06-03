package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username       string `gorm:"uniqueIndex;size:32" json:"username" binding:"required,min=2,max=32"`
	Password       string `json:"password" binding:"required,min=6"`
	Avatar         string `gorm:"size:512" json:"avatar"`
	Bio            string `gorm:"size:500" json:"bio"`
	FollowersCount int    `gorm:"default:0" json:"followersCount"`
	FollowingCount int    `gorm:"default:0" json:"followingCount"`
}

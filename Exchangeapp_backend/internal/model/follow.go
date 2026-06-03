package model

import "time"

type Follow struct {
	ID         uint      `gorm:"primarykey"`
	FollowerID uint      `gorm:"uniqueIndex:idx_follow"`
	FolloweeID uint      `gorm:"uniqueIndex:idx_follow"`
	CreatedAt  time.Time
}

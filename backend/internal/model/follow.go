package model

import "time"

type Follow struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	FollowerID uint      `gorm:"uniqueIndex:idx_follow" json:"followerId"`
	FolloweeID uint      `gorm:"uniqueIndex:idx_follow" json:"followeeId"`
	CreatedAt  time.Time `json:"createdAt"`
}

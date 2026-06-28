package model

import "time"

// PostLike tracks which users have liked which posts.
type PostLike struct {
	ID        uint `gorm:"primarykey"`
	PostID    uint `gorm:"uniqueIndex:idx_post_user"`
	UserID    uint `gorm:"uniqueIndex:idx_post_user"`
	CreatedAt time.Time
}

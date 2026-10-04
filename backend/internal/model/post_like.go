package model

import "time"

// PostLike tracks which users have liked which posts.
type PostLike struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	PostID    uint      `gorm:"uniqueIndex:idx_post_user" json:"postId"`
	UserID    uint      `gorm:"uniqueIndex:idx_post_user" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

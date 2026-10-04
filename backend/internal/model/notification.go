package model

type Notification struct {
	Base
	UserID  uint   `gorm:"index" json:"userId"`
	Type    string `gorm:"size:32" json:"type"`
	Title   string `gorm:"size:256" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	Read    bool   `gorm:"default:false" json:"read"`
}

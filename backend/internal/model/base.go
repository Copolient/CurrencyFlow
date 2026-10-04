package model

import (
	"time"

	"gorm.io/gorm"
)

// Base mirrors gorm.Model but exposes consistent lowerCamelCase JSON fields.
// It keeps the same column names as gorm.Model so AutoMigrate remains compatible.
type Base struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b Base) GetID() uint { return b.ID }

package models

import "time"

type Notification struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`
	Type      string `gorm:"size:50"` // order, system, chat,...
	Title     string
	Content   string
	IsRead    bool `gorm:"index"`
	CreatedAt time.Time
}

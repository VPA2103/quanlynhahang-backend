package models

import "time"

type Notification struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id" gorm:"index;not null"`
	// admin_id hoặc user_id nhận thông báo

	Title   string `json:"title" gorm:"type:varchar(200);not null"`
	Content string `json:"content" gorm:"type:text;not null"`

	Type string `json:"type" gorm:"type:varchar(50);default:'system'"`
	// system | order | contact | payment ...

	IsRead bool `json:"is_read" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at"`
}

package models

import "time"

type PaymentEvent struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index"`
	OrderID   uint
	Status    string // pending, success, failed
	Amount    float64
	CreatedAt time.Time
}

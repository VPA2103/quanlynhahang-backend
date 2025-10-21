package models

import "time"

type HoaDon struct {
	MaHD      string `gorm:"primaryKey;size:10"`
	MaBan     string
	NgayLap   time.Time
	GioLap    string
	TongTien  float64
	TrangThai string
	MaNVOrder string
}

package models

import "time"

type HoaDon struct {
	MaHD      string `gorm:"primaryKey;size:10;autoIncrement"`
	MaBan     string `gorm:"size:10"`
	NgayLap   time.Time
	GioLap    string
	TongTien  float64
	TrangThai string

	MaNVOrder     string   `gorm:"size:10"` // foreign key
	NhanVienOrder NhanVien `gorm:"foreignKey:MaNVOrder;references:MaNV"`
}

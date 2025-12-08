package models

import "time"

type HoaDon struct {
	MaHD      string `gorm:"primaryKey;size:10;autoIncrement"`
	MaBan     string `gorm:"size:10"`
	NgayLap   time.Time
	GioLap    time.Time
	TongTien  float64
	TrangThai uint

	MaNVOrder     uint     `gorm:"size:10"` // foreign key
	NhanVienOrder NhanVien `gorm:"foreignKey:MaNVOrder;references:MaNV"`
}

package models

import "time"

type ThanhToan struct {
	MaThanhToan       string `gorm:"primaryKey;size:10"`
	MaHD              string
	SoTien            float64
	HinhThucThanhToan string
	NgayThanhToan     time.Time
	GioThanhToan      string
}

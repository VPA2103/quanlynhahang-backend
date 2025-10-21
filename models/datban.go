package models

import "time"

type DatBan struct {
	MaDatBan     string `gorm:"primaryKey;size:10"`
	TenKhachHang string
	SDT          string
	GhiChu       string
	MaBan        string
	NgayDatBan   time.Time
	GioDatBan    string
	TrangThai    string
}

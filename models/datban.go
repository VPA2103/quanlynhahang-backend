package models

import "time"

type DatBan struct {
	MaDatBan     uint `gorm:"primaryKey;size:10;autoIncrement"`
	TenKhachHang string
	SDT          string
	GhiChu       string
	MaBanAn      string
	NgayDatBan   time.Time
	GioDatBan    string
	TrangThai    string

	MaNhanVien uint     `gorm:"size:10"` // foreign key
	NhanVien   NhanVien `gorm:"foreignKey:MaNhanVien;references:MaNV"`
}

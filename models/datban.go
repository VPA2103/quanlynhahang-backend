package models

import "time"

type DatBan struct {
	MaDatBan     string `gorm:"primaryKey;size:10"`
	TenKhachHang string
	SDT          string
	GhiChu       string
	MaBanAn      string
	NgayDatBan   time.Time
	GioDatBan    string
	TrangThai    string

	MaNhanVien string   `gorm:"size:10"` // foreign key
	NhanVien   NhanVien `gorm:"foreignKey:MaNhanVien;references:MaNV"`
}

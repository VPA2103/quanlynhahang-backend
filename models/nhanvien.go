package models

import "time"

type NhanVien struct {
	MaNV         string `gorm:"primaryKey;size:10"`
	HoTen        string
	GioiTinh     string
	NgaySinh     time.Time
	SDT          string
	TenDangNhap  string
	DiaChi       string
	NgayVaoLam   time.Time
	Email        string
	MatKhau      string
	AnhNhanVien  string
	LoaiNhanVien string
}

package models

import "time"

type NhanVien struct {
	MaNV         uint `gorm:"primaryKey;size:10;autoIncrement"`
	HoTen        string
	GioiTinh     string
	NgaySinh     time.Time
	SDT          string
	DiaChi       string
	NgayVaoLam   time.Time
	Email        string
	MatKhau      string
	AnhNhanVien  string
	LoaiNhanVien string
}

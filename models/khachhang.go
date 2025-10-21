package models

type KhachHang struct {
	MaKH         string `gorm:"primaryKey;size:10"`
	HoTen        string
	GioiTinh     string
	NgaySinh     string
	DiaChi       string
	Email        string
	MatKhau      string
	AnhKhachHang string
	SDT          string
}

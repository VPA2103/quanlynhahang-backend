package models

type Contact struct {
	MaLienHe     string `gorm:"primaryKey;size:10"`
	TenKhachHang string
	SDT          string
	Email        string
	NoiDung      string
}

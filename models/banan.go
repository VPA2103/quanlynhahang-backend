package models

type BanAn struct {
	MaBan     string `gorm:"primaryKey;size:10;autoIncrementw"`
	TenBan    string
	SoChoNgoi int
	TrangThai string
	AnhBan    string
	Anh_QR    string
}

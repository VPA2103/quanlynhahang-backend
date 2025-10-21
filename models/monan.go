package models

type MonAn struct {
	MaMonAn     string `gorm:"primaryKey;size:10"`
	MaLoaiMonAn string
	TenMonAn    string
	GiaTien     float64
	TrangThai   string
	AnhMonAn    string
}

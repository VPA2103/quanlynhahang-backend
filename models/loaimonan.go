package models

type LoaiMonAn struct {
	MaLoaiMonAn  string `gorm:"primaryKey;size:10"`
	TenLoaiMonAn string
	AnhLoaiMonAn string
}

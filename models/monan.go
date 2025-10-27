package models

type MonAn struct {
	MaMonAn     string `gorm:"primaryKey;size:10;autoIncrement" json:"ma_mon_an"`
	MaLoaiMonAn string `gorm:"size:10"`
	TenMonAn    string
	GiaTien     float64
	TrangThai   string
	AnhMonAn    string
}

package models

type LoaiMonAn struct {
	MaLoaiMonAn  string `gorm:"primaryKey;size:10;autoIncrement" json:"ma_loai_mon_an"`
	TenLoaiMonAn string
	AnhLoaiMonAn string
}

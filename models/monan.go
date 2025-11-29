package models

type MonAn struct {
	MaMonAn     uint    `gorm:"primaryKey;size:10;autoIncrement" json:"ma_mon_an"`
	MaLoaiMonAn uint    `gorm:"size:10" json:"ma_loai_mon_an" form:"ma_loai_mon_an"`
	TenMonAn    string  `json:"ten_mon_an" form:"ten_mon_an"`
	GiaTien     float64 `json:"gia_tien" form:"gia_tien"`
	TrangThai   string  `json:"trang_thai" form:"trang_thai"`
	AnhMonAn    string  `json:"anh_mon_an" form:"anh_mon_an"`
}

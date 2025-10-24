package models

type ChiTietHoaDon struct {
	MaChiTiet string `gorm:"primaryKey;size:10"`
	MaHD      string `gorm:"size:10"`
	MaMonAn   string `gorm:"size:10"`
	SoLuong   int
	DonGia    float64
	ThanhTien float64
	TrangThai string
}

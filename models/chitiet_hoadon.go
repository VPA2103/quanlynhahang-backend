package models

type ChiTietHoaDon struct {
	MaChiTiet string `gorm:"primaryKey;size:10"`
	MaHD      string
	MaMonAn   string
	SoLuong   int
	DonGia    float64
	ThanhTien float64
	TrangThai string
}

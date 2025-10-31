package models

type Images struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	NhanvienID uint   `json:"nhanvien_id"` // Khóa ngoại trỏ tới NhanVien.MaNV
	ImageURL   string `json:"url"`

	// ✅ Dùng pointer + json:"-" để tránh vòng lặp khi marshal JSON
	Nhanvien NhanVien `gorm:"foreignKey:NhanvienID;references:MaNV" json:"-"`
}

package models

type ProductImages struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	NhanvienID uint   `json:"nhanvien_id"` // khóa ngoại
	ImageURL   string `json:"image_url"`

	// dùng pointer để tránh vòng lặp
	Nhanvien *NhanVien `gorm:"foreignKey:NhanvienID;references:MaNV" json:"-"`
}

package models

import "time"

type LienHe struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	HoTen     string    `json:"ho_ten" form:"ho_ten" gorm:"type:varchar(100);not null"`
	Email     string    `json:"email" form:"email" gorm:"type:varchar(150);not null;index"`
	TieuDe    string    `json:"tieu_de" form:"tieu_de" gorm:"type:varchar(200);not null"`
	NoiDung   string    `json:"noi_dung" form:"noi_dung" gorm:"type:text;not null"`
	TrangThai string    `json:"trang_thai" gorm:"type:varchar(50);default:'chưa xử lý'"`
	CreatedAt time.Time `json:"created_at"`
}

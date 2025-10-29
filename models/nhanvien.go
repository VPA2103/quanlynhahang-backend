package models

type NhanVien struct {
	MaNV         uint   `gorm:"primaryKey;autoIncrement" json:"ma_nv" form:"ma_nv"`
	HoTen        string `json:"ho_ten" form:"ho_ten"`
	GioiTinh     string `json:"gioi_tinh" form:"gioi_tinh"`
	NgaySinh     string `json:"ngay_sinh" form:"ngay_sinh"`
	SDT          string `json:"sdt" form:"sdt"`
	DiaChi       string `json:"dia_chi" form:"dia_chi"`
	NgayVaoLam   string `json:"ngay_vao_lam" form:"ngay_vao_lam"`
	Email        string `json:"email" form:"email"`
	MatKhau      string `json:"mat_khau" form:"mat_khau"`
	AnhNhanVien  string `json:"anh_nhan_vien" form:"anh_nhan_vien"`
	LoaiNhanVien string `json:"loai_nhan_vien" form:"loai_nhan_vien"`
}

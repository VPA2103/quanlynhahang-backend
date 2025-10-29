package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"golang.org/x/crypto/bcrypt"
)

// 🧱 Thêm nhân viên
func CreateNhanVien(c *gin.Context) {
	var nv models.NhanVien

	// Đọc dữ liệu từ form
	if err := c.ShouldBindWith(&nv, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Gán NgàyVaoLam = hôm nay nếu chưa có
	if nv.NgayVaoLam == "" {
		nv.NgayVaoLam = time.Now().Format("2006-01-02 15:04:05")
	}

	// ✅ Hash mật khẩu trước khi lưu
	if nv.MatKhau != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(nv.MatKhau), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mã hóa mật khẩu"})
			return
		}
		nv.MatKhau = string(hashedPassword)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mật khẩu không được để trống"})
		return
	}

	// ✅ Lưu vào DB
	if err := config.DB.Create(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo nhân viên thành công",
		"data":    nv,
	})
}

// 📋 Lấy danh sách nhân viên
func GetAllNhanVien(c *gin.Context) {
	var nhanViens []models.NhanVien
	if err := config.DB.Find(&nhanViens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nhanViens)
}

// 🔍 Lấy 1 nhân viên theo ID
func GetNhanVienByID(c *gin.Context) {
	id := c.Param("id")
	var nv models.NhanVien
	if err := config.DB.First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}
	c.JSON(http.StatusOK, nv)
}

// ✏️ Cập nhật nhân viên
func UpdateNhanVien(c *gin.Context) {
	id := c.Param("id")

	// Tìm nhân viên theo ID
	var nv models.NhanVien
	if err := config.DB.First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}

	// Struct tạm để nhận dữ liệu cập nhật (bao gồm cả đổi mật khẩu)
	var req struct {
		HoTen        string `json:"ho_ten" form:"ho_ten"`
		GioiTinh     string `json:"gioi_tinh" form:"gioi_tinh"`
		NgaySinh     string `json:"ngay_sinh" form:"ngay_sinh"`
		SDT          string `json:"sdt" form:"sdt"`
		DiaChi       string `json:"dia_chi" form:"dia_chi"`
		Email        string `json:"email" form:"email"`
		AnhNhanVien  string `json:"anh_nhan_vien" form:"anh_nhan_vien"`
		LoaiNhanVien string `json:"loai_nhan_vien" form:"loai_nhan_vien"`
		OldPassword  string `json:"old_password" form:"old_password"`
		NewPassword  string `json:"new_password" form:"new_password"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// ======================
	// ✅ Nếu có yêu cầu đổi mật khẩu
	// ======================
	if req.NewPassword != "" {
		if req.OldPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng nhập mật khẩu cũ"})
			return
		}

		// Kiểm tra mật khẩu cũ
		if bcrypt.CompareHashAndPassword([]byte(nv.MatKhau), []byte(req.OldPassword)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Mật khẩu cũ không đúng"})
			return
		}

		// Hash mật khẩu mới
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mã hoá mật khẩu mới"})
			return
		}

		nv.MatKhau = string(hashedPassword)
	}

	// ======================
	// ✅ Cập nhật thông tin khác
	// ======================
	if req.HoTen != "" {
		nv.HoTen = req.HoTen
	}
	if req.GioiTinh != "" {
		nv.GioiTinh = req.GioiTinh
	}
	if req.NgaySinh != "" {
		nv.NgaySinh = req.NgaySinh
	}
	if req.SDT != "" {
		nv.SDT = req.SDT
	}
	if req.DiaChi != "" {
		nv.DiaChi = req.DiaChi
	}
	if req.Email != "" {
		nv.Email = req.Email
	}
	if req.AnhNhanVien != "" {
		nv.AnhNhanVien = req.AnhNhanVien
	}
	if req.LoaiNhanVien != "" {
		nv.LoaiNhanVien = req.LoaiNhanVien
	}
	if nv.NgayVaoLam == "" {
		nv.NgayVaoLam = time.Now().Format("2006-01-02 15:04:05")
	}

	// ✅ Lưu thay đổi
	if err := config.DB.Save(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật thông tin nhân viên"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật thông tin nhân viên thành công",
		"data":    nv,
	})
}

// 🗑️ Xóa nhân viên
func DeleteNhanVien(c *gin.Context) {
	id := c.Param("id")
	var nv models.NhanVien
	if err := config.DB.First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}

	if err := config.DB.Delete(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã xóa nhân viên thành công"})
}

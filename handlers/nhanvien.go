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
	var nv models.NhanVien
	if err := config.DB.First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}

	var updateData models.NhanVien
	if err := c.ShouldBind(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&nv).Updates(updateData)
	c.JSON(http.StatusOK, nv)
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

func ChangePasswordNhanVien(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required"`
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu thông tin cần thiết"})
		return
	}

	var nv models.NhanVien
	if err := config.DB.Where("email = ?", req.Email).First(&nv).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}

	// So sánh mật khẩu cũ
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
	if err := config.DB.Save(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật mật khẩu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đổi mật khẩu thành công"})
}

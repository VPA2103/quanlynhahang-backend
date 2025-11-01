package controllers

import (
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"golang.org/x/crypto/bcrypt"
)

// 🧱 Thêm nhân viên
func CreateNhanVien(c *gin.Context) {
	var nv models.NhanVien

	// ✅ Lấy dữ liệu từ form-data
	if err := c.ShouldBind(&nv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu form không hợp lệ: " + err.Error()})
		return
	}

	// ✅ Mặc định ngày vào làm
	if nv.NgayVaoLam == "" {
		nv.NgayVaoLam = time.Now().Format("2006-01-02 15:04:05")
	}

	// ✅ Kiểm tra mật khẩu
	if nv.MatKhau == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mật khẩu không được để trống"})
		return
	}

	// ✅ Hash mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(nv.MatKhau), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mã hóa mật khẩu"})
		return
	}
	nv.MatKhau = string(hashedPassword)

	// ✅ Lưu nhân viên trước để có MaNV (ID)
	if err := config.DB.Create(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo nhân viên: " + err.Error()})
		return
	}

	// ✅ Upload ảnh (nếu có)
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()

			uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
				Folder: "nhanvien",
			})
			if err == nil {
				img := models.Images{
					OwnerID:   nv.MaNV,
					OwnerType: "nhan_vien",
					ImageURL:  uploadResult.SecureURL,
				}
				config.DB.Create(&img)
			}
		}
	}

	// ✅ Preload ảnh khi trả về
	config.DB.Preload("AnhNhanVien").First(&nv, nv.MaNV)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo nhân viên thành công",
		"data":    nv,
	})
}

// 📋 Lấy danh sách nhân viên
func GetAllNhanVien(c *gin.Context) {
	var nhanViens []models.NhanVien
	if err := config.DB.Preload("AnhNhanVien").Find(&nhanViens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nhanViens)
}

// 🔍 Lấy 1 nhân viên theo ID
func GetNhanVienByID(c *gin.Context) {
	id := c.Param("id")
	var nv models.NhanVien
	if err := config.DB.Preload("AnhNhanVien").Find(&nv, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nv)
}

// ✏️ Cập nhật nhân viên
func UpdateNhanVien(c *gin.Context) {
	id := c.Param("id")
	var nv models.NhanVien

	// ✅ Tìm nhân viên theo ID
	if err := config.DB.Preload("AnhNhanVien").First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}

	// ✅ Bind dữ liệu từ form
	var updatedData models.NhanVien
	if err := c.ShouldBind(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu form không hợp lệ: " + err.Error()})
		return
	}

	// ✅ Cập nhật các trường cơ bản
	nv.HoTen = updatedData.HoTen
	nv.GioiTinh = updatedData.GioiTinh
	nv.NgaySinh = updatedData.NgaySinh
	nv.SDT = updatedData.SDT
	nv.DiaChi = updatedData.DiaChi
	nv.Email = updatedData.Email
	nv.LoaiNhanVien = updatedData.LoaiNhanVien

	// ✅ Nếu có thay đổi mật khẩu
	if updatedData.MatKhau != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(updatedData.MatKhau), bcrypt.DefaultCost)
		nv.MatKhau = string(hashedPassword)
	}

	// ✅ Upload ảnh mới (nếu có)
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, _ := file.Open()
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{Folder: "nhanvien"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload ảnh thất bại: " + err.Error()})
			return
		}

		// 🔹 Xóa ảnh cũ trong DB (nếu có)
		config.DB.Where("owner_id = ? AND owner_type = ?", nv.MaNV, "nhan_vien").Delete(&models.Images{})

		// 🔹 Lưu ảnh mới
		newImg := models.Images{
			OwnerID:   nv.MaNV,
			OwnerType: "nhan_vien",
			ImageURL:  uploadResult.SecureURL,
		}
		config.DB.Create(&newImg)
	}

	// ✅ Lưu thay đổi nhân viên
	if err := config.DB.Save(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật nhân viên: " + err.Error()})
		return
	}

	// ✅ Lấy lại nhân viên kèm ảnh mới
	config.DB.Preload("AnhNhanVien").First(&nv, nv.MaNV)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật nhân viên thành công",
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

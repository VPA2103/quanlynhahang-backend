package controllers

import (
	"fmt"
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

	if err := c.ShouldBind(&nv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu form không hợp lệ: " + err.Error()})
		return
	}

	if nv.NgayVaoLam == "" {
		nv.NgayVaoLam = time.Now().Format("2006-01-02 15:04:05")
	}

	if nv.MatKhau == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mật khẩu không được để trống"})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(nv.MatKhau), bcrypt.DefaultCost)
	nv.MatKhau = string(hashedPassword)

	if err := config.DB.Create(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo nhân viên: " + err.Error()})
		return
	}

	// Upload ảnh (nếu có)
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, _ := file.Open()
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{Folder: "nhanvien"})
		if err == nil {
			img := models.Images{
				NhanvienID: nv.MaNV,
				ImageURL:   uploadResult.SecureURL,
			}
			config.DB.Create(&img)
		}
	}

	// ✅ Lấy lại nhân viên kèm ảnh
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
	if err := config.DB.First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}
	c.JSON(http.StatusOK, nv)
}

// ✏️ Cập nhật nhân viên
func UpdateNhanVien(c *gin.Context) {
	id := c.Param("id")

	// 🔹 Lấy role từ token (JWT)
	roleVal, _ := c.Get("role")
	role := fmt.Sprintf("%v", roleVal)

	// 🔹 Lấy user ID từ token (để giới hạn quyền)
	userIDVal, _ := c.Get("user_id")
	userID := uint(0)
	if uid, ok := userIDVal.(float64); ok {
		userID = uint(uid)
	}

	// 🔹 Tìm nhân viên theo ID
	var nv models.NhanVien
	if err := config.DB.Preload("AnhNhanVien").First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}

	// 🔐 Giới hạn: nếu không phải admin thì chỉ được sửa chính mình
	if role != "admin" && userID != nv.MaNV {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền chỉnh sửa người khác"})
		return
	}

	// 🔹 Bind dữ liệu form
	var req struct {
		HoTen        string `json:"ho_ten" form:"ho_ten"`
		GioiTinh     string `json:"gioi_tinh" form:"gioi_tinh"`
		NgaySinh     string `json:"ngay_sinh" form:"ngay_sinh"`
		SDT          string `json:"sdt" form:"sdt"`
		DiaChi       string `json:"dia_chi" form:"dia_chi"`
		Email        string `json:"email" form:"email"`
		LoaiNhanVien string `json:"loai_nhan_vien" form:"loai_nhan_vien"`
		OldPassword  string `json:"old_password" form:"old_password"`
		NewPassword  string `json:"new_password" form:"new_password"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	// ==========================
	// ✅ Xử lý upload ảnh mới
	// ==========================
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mở file ảnh"})
			return
		}
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{Folder: "nhanvien"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload ảnh thất bại: " + err.Error()})
			return
		}

		if len(nv.AnhNhanVien) > 0 {
			nv.AnhNhanVien[0].ImageURL = uploadResult.SecureURL
			config.DB.Save(&nv.AnhNhanVien[0])
		} else {
			newImg := models.Images{
				NhanvienID: nv.MaNV,
				ImageURL:   uploadResult.SecureURL,
			}
			config.DB.Create(&newImg)
			nv.AnhNhanVien = append(nv.AnhNhanVien, newImg)
		}
	}

	// ==========================
	// ✅ Xử lý đổi mật khẩu
	// ==========================
	if req.NewPassword != "" {
		if role == "admin" {
			// 🔓 Admin đổi mật khẩu mà không cần mật khẩu cũ
			hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
			nv.MatKhau = string(hashed)
		} else {
			// 🧱 Nhân viên thường phải nhập mật khẩu cũ
			if req.OldPassword == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng nhập mật khẩu cũ"})
				return
			}
			if bcrypt.CompareHashAndPassword([]byte(nv.MatKhau), []byte(req.OldPassword)) != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Mật khẩu cũ không đúng"})
				return
			}
			hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
			nv.MatKhau = string(hashed)
		}
	}

	// ==========================
	// ✅ Cập nhật thông tin khác
	// ==========================
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
	if req.LoaiNhanVien != "" && role == "admin" {
		// 🧱 Chỉ admin mới được thay đổi loại nhân viên
		nv.LoaiNhanVien = req.LoaiNhanVien
	}

	// ==========================
	// ✅ Lưu thay đổi
	// ==========================
	if err := config.DB.Save(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu thông tin nhân viên"})
		return
	}

	config.DB.Preload("AnhNhanVien").First(&nv, id)
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
